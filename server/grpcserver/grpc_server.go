package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/w-h-a/pkg/server"
	"github.com/w-h-a/pkg/server/grpcserver/handlers"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/marshalutils"
	"github.com/w-h-a/pkg/utils/metadatautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	defaultContentType = "application/grpc+proto"
)

type grpcServer struct {
	options  server.ServerOptions
	server   *grpc.Server
	handlers map[string]*grpcHandler
	started  bool
	exit     chan chan error
	mtx      sync.RWMutex
	wg       *sync.WaitGroup
}

func (s *grpcServer) Options() server.ServerOptions {
	return s.options
}

func (s *grpcServer) NewHandler(c interface{}, opts ...server.HandlerOption) server.Handler {
	return NewHandler(c, opts...)
}

func (s *grpcServer) Handle(h server.Handler) error {
	handler, ok := h.(*grpcHandler)
	if !ok {
		return fmt.Errorf("invalid handler: expected *grpcHandler")
	}

	if len(handler.methods) == 0 {
		return fmt.Errorf("invalid handler: no methods")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.handlers[h.Name()]; ok {
		return fmt.Errorf("handler %+v is already registered", handler)
	}

	s.handlers[h.Name()] = handler

	return nil
}

func (s *grpcServer) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Infof("grpc server received signal %s", <-ch)

	return s.stop()
}

func (s *grpcServer) String() string {
	return "grpc"
}

func (s *grpcServer) start() error {
	s.mtx.RLock()
	if s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	if err := handlers.RegisterHealthHandler(
		s,
		handlers.NewHealthHandler(
			fmt.Sprintf("%s.%s:%s %s", s.options.Name, s.options.Namespace, s.options.Version, s.options.Id),
		),
	); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", s.options.Address)
	if err != nil {
		return err
	}

	// make sure the address is right
	s.mtx.Lock()
	s.options.Address = listener.Addr().String()
	s.mtx.Unlock()

	log.Infof("grpc server is listening on %s", listener.Addr().String())

	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("grpc server failed to start: %v", err)
		}
	}()

	go func() {
		tick := time.NewTicker(time.Second * 30)

		var ch chan error

	Loop:
		for {
			select {
			case <-tick.C:
				// TODO: decide if we need this for some kind of heart beat
			case ch = <-s.exit:
				tick.Stop()
				break Loop
			}
		}

		wait := make(chan struct{})

		go func() {
			defer close(wait)
			s.wg.Wait()
		}()

		select {
		case <-wait:
		case <-time.After(30 * time.Second):
		}

		shutdown := make(chan struct{})

		go func() {
			defer close(shutdown)
			s.server.GracefulStop()
		}()

		select {
		case <-shutdown:
		case <-time.After(30 * time.Second):
			s.server.Stop()
		}

		// signal that we've finished stopping the grpc server
		ch <- nil
	}()

	s.mtx.Lock()
	s.started = true
	s.mtx.Unlock()

	return nil
}

func (s *grpcServer) stop() error {
	s.mtx.RLock()
	if !s.started {
		s.mtx.RUnlock()
		return nil
	}
	s.mtx.RUnlock()

	ch := make(chan error)

	// signal start loop
	s.exit <- ch

	// wait for errors
	err := <-ch

	s.mtx.Lock()
	s.started = false
	s.mtx.Unlock()

	return err
}

func (s *grpcServer) handle(_ interface{}, stream grpc.ServerStream) error {
	s.wg.Add(1)
	defer s.wg.Done()

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

	contentType := defaultContentType
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

	method := handler.methods[methodName]

	if method == nil {
		return status.New(codes.Unimplemented, fmt.Sprintf("unknown method %s.%s", handlerName, methodName)).Err()
	}

	if method.stream {
		return s.processStream(stream, handler, method, contentType, ctx)
	}

	return s.processRequest(stream, handler, method, contentType, ctx)
}

func (s *grpcServer) processRequest(stream grpc.ServerStream, handler *grpcHandler, method *grpcMethod, contentType string, ctx context.Context) error {
	req := reflect.New(method.reqType.Elem())

	rsp := reflect.New(method.rspType.Elem())

	if err := stream.RecvMsg(req.Interface()); err != nil {
		return err
	}

	// this is necessary in addition to the init toward the
	// bottom to get grpc to assume the right content type
	marshaler, err := s.newMarshaler(contentType)
	if err != nil {
		return errorutils.InternalServerError("server", err.Error())
	}

	bytes, err := marshaler.Marshal(req.Interface())
	if err != nil {
		return err
	}

	fn := func(ctx context.Context, request server.Request, response interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Infof("panic recovered: %v", r)
				log.Info(string(debug.Stack()))
				err = errorutils.InternalServerError("server", "panic recovered: %v", r)
			}
		}()

		args := []reflect.Value{
			handler.receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(request.Unmarshaled()),
			reflect.ValueOf(response),
		}

		vals := method.value.Call(args)

		if e := vals[0].Interface(); e != nil {
			err = e.(error)
		}

		return
	}

	for i := len(s.options.HandlerWrappers); i > 0; i-- {
		fn = s.options.HandlerWrappers[i-1](fn)
	}

	statusCode := codes.OK
	statusDesc := ""

	if err := fn(
		ctx,
		NewRequest(
			server.RequestWithNamespace(s.options.Namespace),
			server.RequestWithName(s.options.Name),
			server.RequestWithMethod(fmt.Sprintf("%s.%s", handler.name, method.name)),
			server.RequestWithContentType(contentType),
			server.RequestWithUnmarshaledRequest(req.Interface()),
			server.RequestWithMarshaledRequest(bytes),
		),
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

func (s *grpcServer) processStream(stream grpc.ServerStream, handler *grpcHandler, method *grpcMethod, contentType string, ctx context.Context) error {
	fn := func(ctx context.Context, request server.Request, stream interface{}) (err error) {
		args := []reflect.Value{
			handler.receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(stream),
		}

		vals := method.value.Call(args)

		if e := vals[0].Interface(); e != nil {
			err = e.(error)
		}

		return
	}

	for i := len(s.options.HandlerWrappers); i > 0; i-- {
		fn = s.options.HandlerWrappers[i-1](fn)
	}

	statusCode := codes.OK
	statusDesc := ""

	r := NewRequest(
		server.RequestWithNamespace(s.options.Namespace),
		server.RequestWithName(s.options.Name),
		server.RequestWithMethod(fmt.Sprintf("%s.%s", handler.name, method.name)),
		server.RequestWithContentType(contentType),
		server.RequestWithStream(),
	)

	if err := fn(
		ctx,
		r,
		&grpcStream{
			request: r,
			stream:  stream,
		},
	); err != nil {
		statusCode = ToErrorCode(err)
		statusDesc = err.Error()
		return status.New(statusCode, statusDesc).Err()
	}

	return status.New(statusCode, statusDesc).Err()
}

func (s *grpcServer) newMarshaler(contentType string) (encoding.Codec, error) {
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

func NewServer(opts ...server.ServerOption) server.Server {
	options := server.NewServerOptions(opts...)

	s := &grpcServer{
		options:  options,
		handlers: map[string]*grpcHandler{},
		exit:     make(chan chan error),
		mtx:      sync.RWMutex{},
		wg:       &sync.WaitGroup{},
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnknownServiceHandler(s.handle),
	}

	s.server = grpc.NewServer(grpcOptions...)

	return s
}
