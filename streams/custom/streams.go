package custom

import (
	"context"
	"encoding/json"
	"time"

	"github.com/w-h-a/pkg/client"
	pb "github.com/w-h-a/pkg/proto/streams"
	"github.com/w-h-a/pkg/streams"
	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	OrderCreated   = "order:created"
	OrderCancelled = "order:cancelled"
	OrderExpired   = "order:expired"
	PaymentCreated = "payment:created"
)

type customStream struct {
	options streams.StreamsOptions
	streams client.Client
}

func (s *customStream) Options() streams.StreamsOptions {
	return s.options
}

func (s *customStream) Subscribe(id string, opts ...streams.SubscribeOption) error {
	options := streams.NewSubscribeOptions(opts...)

	req := s.streams.NewRequest(
		client.RequestWithNamespace("wha-platform-resource"),
		client.RequestWithPort(8081),
		client.RequestWithName("streams"),
		client.RequestWithMethod("Stream.Subscribe"),
		client.RequestWithUnmarshaledRequest(
			&pb.SubscribeRequest{
				Id:         id,
				Group:      options.Group,
				Topic:      options.Topic,
				AckWait:    options.AckWait.Nanoseconds(),
				RetryLimit: int64(options.RetryLimit),
				Offset:     options.Offset.Unix(),
			},
		),
	)

	rsp := &pb.SubscribeResponse{}

	if err := s.streams.Call(context.Background(), req, rsp); err != nil {
		return err
	}

	return nil
}

func (s *customStream) Unsubscribe(id string) error {
	req := s.streams.NewRequest(
		client.RequestWithNamespace("wha-platform-resource"),
		client.RequestWithPort(8081),
		client.RequestWithName("streams"),
		client.RequestWithMethod("Stream.Unsubscribe"),
		client.RequestWithUnmarshaledRequest(
			&pb.UnsubscribeRequest{
				Id: id,
			},
		),
	)

	rsp := &pb.UnsubscribeResponse{}

	if err := s.streams.Call(context.Background(), req, rsp); err != nil {
		return err
	}

	return nil
}

func (s *customStream) Consume(id string) (streams.Subscriber, error) {
	pbReq := &pb.ConsumeRequest{
		Id: id,
	}

	req := s.streams.NewRequest(
		client.RequestWithNamespace("wha-platform-resource"),
		client.RequestWithPort(8081),
		client.RequestWithName("streams"),
		client.RequestWithMethod("Stream.Consume"),
		client.RequestWithUnmarshaledRequest(pbReq),
	)

	stream, err := s.streams.Stream(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if err := stream.Send(pbReq); err != nil {
		return nil, err
	}

	sub := NewSubscriber()

	go func() {
		for {
			var ev pb.Event

			if err := stream.Recv(&ev); err != nil {
				log.Errorf("failed to receive event from stream: %v", err)
				sub.Close()
				stream.Close()
				return
			}

			event := &streams.Event{
				Id:        ev.Id,
				Topic:     ev.Topic,
				Payload:   ev.Payload,
				Timestamp: time.Unix(ev.Timestamp, 0),
				Metadata:  ev.Metadata,
			}

			cpy := *event

			cpy.SetAck(Ack(stream, sub, cpy))
			cpy.SetNack(Nack(stream, sub, cpy))

			sub.Channel() <- cpy
		}
	}()

	return sub, nil
}

func (s *customStream) Produce(topic string, data interface{}, opts ...streams.ProduceOption) error {
	options := streams.NewProduceOptions(opts...)

	var bytes []byte

	if p, ok := data.([]byte); ok {
		bytes = p
	} else {
		p, err := json.Marshal(data)
		if err != nil {
			return streams.ErrEncodingData
		}
		bytes = p
	}

	req := s.streams.NewRequest(
		client.RequestWithNamespace("wha-platform-resource"),
		client.RequestWithPort(8081),
		client.RequestWithName("streams"),
		client.RequestWithMethod("Stream.Produce"),
		client.RequestWithUnmarshaledRequest(
			&pb.ProduceRequest{
				Topic:    topic,
				Payload:  bytes,
				Metadata: options.Metadata,
			},
		),
	)

	rsp := &pb.ProduceResponse{}

	if err := s.streams.Call(options.Context, req, rsp); err != nil {
		return err
	}

	return nil
}

func (s *customStream) String() string {
	return "custom"
}

func NewStreams(c client.Client) streams.Streams {
	o := streams.NewStreamsOptions()
	return &customStream{o, c}
}
