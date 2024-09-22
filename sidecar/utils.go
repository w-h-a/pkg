package sidecar

import (
	"encoding/json"

	pb "github.com/w-h-a/pkg/proto/sidecar"
	"google.golang.org/protobuf/types/known/anypb"
)

func SerializeEvent(event *Event) (*pb.Event, error) {
	bs, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}

	return &pb.Event{
		EventName: event.EventName,
		Data: &anypb.Any{
			Value: bs,
		},
	}, nil
}
