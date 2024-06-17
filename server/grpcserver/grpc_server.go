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
	"github.com/w-h-a/pkg/server/grpcserver/controllers"
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
	options     server.ServerOptions
	mtx         sync.RWMutex
	wg          *sync.WaitGroup
	controllers map[string]*grpcController
	server      *grpc.Server
	started     bool
	exit        chan chan error
}

func (s *grpcServer) Options() server.ServerOptions {
	return s.options
}

func (s *grpcServer) NewController(c interface{}, opts ...server.ControllerOption) server.Controller {
	return NewController(c, opts...)
}

func (s *grpcServer) RegisterController(c server.Controller) error {
	controller, ok := c.(*grpcController)
	if !ok {
		return fmt.Errorf("invalid controller: expected *grpcController")
	}

	if len(controller.handlers) == 0 {
		return fmt.Errorf("invalid controller: no handler functions")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.controllers[c.Name()]; ok {
		return fmt.Errorf("controller %+v is already registered", controller)
	}

	s.controllers[c.Name()] = controller

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

	if err := controllers.RegisterHealthController(
		s,
		controllers.NewHealthController(
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

	// TODO: connect to broker if we have subscribers

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

		if s.wg != nil {
			s.wg.Wait()
		}

		s.server.GracefulStop()

		// signal that we've finished stopping the grpc server
		ch <- nil

		// TODO: disconnect from broker if we're connected
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

	controllerName, handlerName, err := ToControllerHandler(grpcFormattedMethod)
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
	controller := s.controllers[controllerName]
	s.mtx.Unlock()

	if controller == nil {
		return status.New(codes.Unimplemented, fmt.Sprintf("unknown controller %s", controllerName)).Err()
	}

	handler := controller.handlers[handlerName]

	if handler == nil {
		return status.New(codes.Unimplemented, fmt.Sprintf("unknown method %s.%s", controllerName, handlerName)).Err()
	}

	return s.processRequest(stream, controller, handler, contentType, ctx)
}

func (s *grpcServer) processRequest(stream grpc.ServerStream, controller *grpcController, handler *grpcSync, contentType string, ctx context.Context) error {
	req := reflect.New(handler.reqType.Elem())

	rsp := reflect.New(handler.rspType.Elem())

	if err := stream.RecvMsg(req.Interface()); err != nil {
		return err
	}

	marshaler, err := s.newMarshaler(contentType)
	if err != nil {
		return errorutils.InternalServerError("server", err.Error())
	}

	b, err := marshaler.Marshal(req.Interface())
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
			controller.receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(request.Unmarshaled()),
			reflect.ValueOf(response),
		}

		results := handler.method.Call(args)

		if e := results[0].Interface(); e != nil {
			err = e.(error)
		}

		return err
	}

	for i := len(s.options.ControllerWrappers); i > 0; i-- {
		fn = s.options.ControllerWrappers[i-1](fn)
	}

	statusCode := codes.OK
	statusDesc := ""

	if err := fn(
		ctx,
		NewRequest(
			server.RequestWithNamespace(s.options.Namespace),
			server.RequestWithName(s.options.Name),
			server.RequestWithMethod(fmt.Sprintf("%s.%s", controller.name, handler.name)),
			server.RequestWithContentType(contentType),
			server.RequestWithUnmarshaledRequest(req.Interface()),
			server.RequestWithMarshaledRequest(b),
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

func (s *grpcServer) newMarshaler(contentType string) (encoding.Codec, error) {
	marshaler, ok := marshalutils.DefaultMarshalers[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return marshaler, nil
}

// func init() {
// 	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/json"])
// 	encoding.RegisterCodec(marshalutils.DefaultMarshalers["application/proto"])
// }

func NewServer(opts ...server.ServerOption) server.Server {
	options := server.NewServerOptions(opts...)

	s := &grpcServer{
		options:     options,
		mtx:         sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		controllers: map[string]*grpcController{},
		exit:        make(chan chan error),
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnknownServiceHandler(s.handle),
	}

	s.server = grpc.NewServer(grpcOptions...)

	return s
}
