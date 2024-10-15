package sidecar

import (
	pb "github.com/w-h-a/pkg/proto/sidecar"
)

func SerializeEvent(event *Event) (*pb.Event, error) {
	return &pb.Event{
		EventName: event.EventName,
		Payload: &pb.Payload{
			Metadata: event.Payload.Metadata,
			Data:     event.Payload.Data,
		},
	}, nil
}
