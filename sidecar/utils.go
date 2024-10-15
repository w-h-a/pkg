package sidecar

import (
	"encoding/json"

	pb "github.com/w-h-a/pkg/proto/sidecar"
	"google.golang.org/protobuf/types/known/anypb"
)

func SerializeEvent(event *Event) (*pb.Event, error) {
	bs, _ := json.Marshal(event.Payload.Data)

	return &pb.Event{
		EventName: event.EventName,
		Payload: &pb.Payload{
			Metadata: event.Payload.Metadata,
			Data: &anypb.Any{
				Value: bs,
			},
		},
	}, nil
}
