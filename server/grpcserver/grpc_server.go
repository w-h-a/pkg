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

	"github.com/w-h-a/pkg/broker"
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
	subscribers map[*grpcSubscriber]broker.Subscriber
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

func (s *grpcServer) NewSubscriber(t string, sub interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return NewSubscriber(t, sub, opts...)
}

func (s *grpcServer) RegisterSubscriber(sub server.Subscriber) error {
	subscriber, ok := sub.(*grpcSubscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *grpcSubscriber")
	}

	if len(subscriber.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.subscribers[subscriber]; ok {
		return fmt.Errorf("subscriber %v is already registered", subscriber)
	}

	s.subscribers[subscriber] = nil

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

	if len(s.subscribers) > 0 && s.options.Broker != nil {
		if err := s.options.Broker.Connect(); err != nil {
			log.Errorf("grpc server failed to connect to broker: %v", err)
			return err
		}

		log.Infof("grpc server is connected to [%s] broker at %s", s.options.Broker.String(), s.options.Broker.Address())

		if err := s.subscribe(); err != nil {
			log.Errorf("grpc server failed to subscribe subscribers: %v", err)
			return err
		}
	}

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

		if s.options.Broker != nil && len(s.subscribers) != 0 {
			if err := s.unsubscribe(); err != nil {
				log.Errorf("grpc server failed to unsubscribe subscribers: %v", err)
			}

			if err := s.options.Broker.Disconnect(); err != nil {
				log.Errorf("grpc server failed to disconnect from the broker: %v", err)
			}
		}
	}()

	s.mtx.Lock()
	s.started = true
	s.mtx.Unlock()

	return nil
}

func (s *grpcServer) subscribe() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for grpc := range s.subscribers {
		handler := s.createHandler(grpc, s.options.SubscriberWrappers)

		opts := []broker.SubscribeOption{}

		if !grpc.Options().AutoAck {
			opts = append(opts, broker.SubscribeWithoutAutoAck())
		}

		if len(grpc.Options().Queue) > 0 {
			opts = append(opts, broker.SubscribeWithQueue(grpc.Options().Queue))
		}

		log.Infof("grpc server subscribing to topic: %s", grpc.Topic())

		sub, err := s.options.Broker.Subscribe(grpc.Topic(), handler, opts...)
		if err != nil {
			return err
		}

		s.subscribers[grpc] = sub
	}

	return nil
}

func (s *grpcServer) unsubscribe() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for grpc, sub := range s.subscribers {
		log.Infof("grpc server unsubscribing from topic: %s", sub.Topic())

		if err := sub.Unsubscribe(); err != nil {
			return err
		}

		s.subscribers[grpc] = nil
	}

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

	// seems like this was necessary, in addition to the init toward the
	// bottom of this file, to get grpc to assume the right content type
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
			controller.receiver,
			reflect.ValueOf(ctx),
			reflect.ValueOf(request.Unmarshaled()),
			reflect.ValueOf(response),
		}

		vals := handler.method.Call(args)

		if e := vals[0].Interface(); e != nil {
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

func (s *grpcServer) createHandler(subscriber *grpcSubscriber, subWrappers []server.SubscriberWrapper) broker.Handler {
	return func(pub broker.Publication) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Infof("panic recovered: %v", r)
				log.Info(string(debug.Stack()))
				err = errorutils.InternalServerError("server", "panic recovered: %v", r)
			}
		}()

		message := pub.Message()
		if message.Header == nil {
			message.Header = map[string]string{}
		}

		md := metadatautils.Metadata{}
		for k, v := range message.Header {
			md[k] = v
		}

		contentType := defaultContentType
		if ct, ok := md["content-type"]; ok {
			contentType = ct
		}

		ctx := metadatautils.NewContext(context.Background(), md)

		var marshaler marshalutils.Marshaler
		marshaler, err = s.newMarshaler(contentType)
		if err != nil {
			return
		}

		results := make(chan error, len(subscriber.handlers))

		for i := 0; i < len(subscriber.handlers); i++ {
			handler := subscriber.handlers[i]

			payload := reflect.New(handler.payloadType.Elem())

			err = marshaler.Unmarshal(message.Body, payload.Interface())
			if err != nil {
				return
			}

			fn := func(ctx context.Context, pub server.Publication) (err error) {
				args := []reflect.Value{
					subscriber.receiver,
					reflect.ValueOf(ctx),
					reflect.ValueOf(pub.Unmarshaled()),
				}

				vals := handler.method.Call(args)

				if e := vals[0].Interface(); e != nil {
					err = e.(error)
				}

				return err
			}

			for i := len(subWrappers); i > 0; i-- {
				fn = subWrappers[i-1](fn)
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				e := fn(
					ctx,
					NewPublication(
						server.PublicationWithTopic(subscriber.topic),
						server.PublicationWithContentType(contentType),
						server.PublicationWithUnmarshaledPayload(payload.Interface()),
					),
				)
				results <- e
			}()
		}

		errors := []string{}

		for i := 0; i < len(subscriber.handlers); i++ {
			if e := <-results; e != nil {
				errors = append(errors, e.Error())
			}
		}

		if len(errors) > 0 {
			err = fmt.Errorf("subscriber error: %s", strings.Join(errors, "\n"))
			log.Errorf("subscriber errors: %v", err)
		}

		return
	}
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
		options:     options,
		mtx:         sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		controllers: map[string]*grpcController{},
		subscribers: map[*grpcSubscriber]broker.Subscriber{},
		exit:        make(chan chan error),
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnknownServiceHandler(s.handle),
	}

	s.server = grpc.NewServer(grpcOptions...)

	return s
}
