package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/w-h-a/pkg/serverv2"
	"github.com/w-h-a/pkg/serverv2/grpc/handlers"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/trace"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/marshalutils"
	"github.com/w-h-a/pkg/utils/metadatautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	options  serverv2.ServerOptions
	handlers map[string]*Handler
	started  bool
	mtx      sync.RWMutex
	errCh    chan error
	exit     chan struct{}
}

func (s *server) Options() serverv2.ServerOptions {
	return s.options
}

func (s *server) Handle(h interface{}) error {
	handler, ok := h.(*Handler)
	if !ok {
		return fmt.Errorf("invalid handler: expected *grpc.Handler")
	}

	if len(handler.Methods) == 0 {
		return fmt.Errorf("invalid handler: no exported methods were found")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.handlers[handler.Name]; ok {
		return fmt.Errorf("handler %#+v is already registered", handler)
	}

	s.handlers[handler.Name] = handler

	return nil
}

func (s *server) Start() error {
	if err := s.Run(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Infof("grpc server received signal %s", <-ch)

	return s.Stop()
}

func (s *server) Run() error {
	s.mtx.RLock()
	if s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	if len(s.options.Tracer) > 0 {
		// init trace exporters
		switch s.options.Tracer {
		case "memory":
			tracer := trace.GetTracer()
			if tracer == nil {
				log.Fatalf("failed to init memory trace exporter: memory tracer was not set")
			}
			grpcTrace := handlers.NewTraceHandler(tracer)
			s.Handle(NewHandler(grpcTrace))
		default:
			log.Warnf("tracer %s is not supported", s.options.Tracer)
		}
	}

	// TODO: tls
	listener, err := net.Listen("tcp", s.options.Address)
	if err != nil {
		return err
	}

	s.mtx.Lock()
	s.options.Address = listener.Addr().String()
	s.mtx.Unlock()

	log.Infof("grpc server is listening on %s", s.options.Address)

	grpcServer := grpc.NewServer(grpc.UnknownServiceHandler(s.handle))

	go func() {
		s.errCh <- grpcServer.Serve(listener)
	}()

	go func() {
		<-s.exit

		var err error

		shutdown := make(chan struct{})

		go func() {
			defer close(shutdown)
			grpcServer.GracefulStop()
		}()

		select {
		case <-shutdown:
		case <-time.After(10 * time.Second):
			grpcServer.Stop()
		}

		s.errCh <- err
	}()

	s.mtx.Lock()
	s.started = true
	s.mtx.Unlock()

	return nil
}

func (s *server) Stop() error {
	s.mtx.RLock()
	if !s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	close(s.exit)

	var err error

	err = <-s.errCh
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, grpc.ErrServerStopped) {
		err = nil
	}

	s.mtx.Lock()
	s.started = false
	s.mtx.Unlock()

	return err
}

func (s *server) String() string {
	return "grpc"
}

func (s *server) handle(_ interface{}, stream grpc.ServerStream) error {
	grpcFormattedMethod, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Errorf(codes.Internal, "method is not present in context")
	}

	handlerName, methodName, err := ToHandlerMethod(grpcFormattedMethod)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error()).Err()
	}

	grpcMetadata, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		grpcMetadata = metadata.MD{}
	}

	md := metadatautils.Metadata{}
	for k, v := range grpcMetadata {
		md[k] = strings.Join(v, ", ")
	}

	contentType := "application/grpc+proto"
	if ct, ok := md["content-type"]; ok {
		contentType = ct
	}

	timeout := md["timeout"]
	delete(md, "timeout")

	ctx := metadatautils.NewContext(stream.Context(), md)

	if len(timeout) > 0 {
		n, err := strconv.ParseUint(timeout, 10, 64)
		if err == nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Duration(n))
			defer cancel()
		}
	}

	s.mtx.Lock()
	handler := s.handlers[handlerName]
	s.mtx.Unlock()

	if handler == nil {
		return status.New(codes.Unimplemented, fmt.Sprintf("unknown handler %s", handlerName)).Err()
	}

	method := handler.Methods[methodName]

	if method == nil {
		return status.New(codes.Unimplemented, fmt.Sprintf("unknown method %s.%s", handlerName, methodName)).Err()
	}

	// TODO: stream
	// if method.stream {
	// 	return s.processStream(stream, handler, method, contentType, ctx)
	// }

	return s.processRequest(stream, handler, method, contentType, ctx)
}

func (s *server) processRequest(stream grpc.ServerStream, handler *Handler, method *Method, contentType string, ctx context.Context) error {
	req := reflect.New(method.ReqType.Elem())

	rsp := reflect.New(method.RspType.Elem())

	if err := stream.RecvMsg(req.Interface()); err != nil {
		return err
	}

	// this is necessary in addition to the init toward the
	// bottom to get grpc to assume the right content type
	marshaler, err := s.newMarshaler(contentType)
	if err != nil {
		return errorutils.InternalServerError("server", err.Error())
	} else if _, err := marshaler.Marshal(req.Interface()); err != nil {
		return errorutils.InternalServerError("server", err.Error())
	}

	fun := func(ctx context.Context, request interface{}, response interface{}) (err error) {
		args := []reflect.Value{
			handler.Receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(request),
			reflect.ValueOf(response),
		}

		vals := method.Value.Call(args)

		if e := vals[0].Interface(); e != nil {
			err = e.(error)
		}

		return
	}

	if ms, ok := GetMiddlewaresFromContext(s.options.Context); ok && ms != nil {
		for i := len(ms); i > 0; i-- {
			fun = ms[i-1](fun)
		}
	}

	statusCode := codes.OK
	statusDesc := ""

	if err := fun(
		ctx,
		req.Interface(),
		rsp.Interface(),
	); err != nil {
		statusCode = ToErrorCode(err)
		statusDesc = err.Error()
		return status.New(statusCode, statusDesc).Err()
	}

	if err := stream.SendMsg(rsp.Interface()); err != nil {
		return err
	}

	return status.New(statusCode, statusDesc).Err()
}

func (s *server) newMarshaler(contentType string) (encoding.Codec, error) {
	marshaler, ok := marshalutils.DefaultMarshalers[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return marshaler, nil
}

func init() {
	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/json"])
	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/proto"])
}

func NewServer(opts ...serverv2.ServerOption) serverv2.Server {
	options := serverv2.NewServerOptions(opts...)

	s := &server{
		options:  options,
		handlers: map[string]*Handler{},
		mtx:      sync.RWMutex{},
		errCh:    make(chan error),
		exit:     make(chan struct{}),
	}

	return s
}
