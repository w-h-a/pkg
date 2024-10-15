package sidecar

import (
	"encoding/json"

	pb "github.com/w-h-a/pkg/proto/sidecar"
	"google.golang.org/protobuf/types/known/anypb"
)

func SerializeEvent(event *Event) (*pb.Event, error) {
	kvs := []*pb.KeyVal{}

	for k, v := range event.Payload {
		bytes, _ := json.Marshal(v)
		kv := &pb.KeyVal{Key: k, Value: &anypb.Any{Value: bytes}}
		kvs = append(kvs, kv)
	}

	return &pb.Event{
		EventName: event.EventName,
		Payload:   kvs,
	}, nil
}
