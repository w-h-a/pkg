package sidecar

import (
	"encoding/json"

	pb "github.com/w-h-a/pkg/proto/sidecar"
)

func SerializeEvent(event *Event) (*pb.Event, error) {
	bs, _ := json.Marshal(event.Payload)

	return &pb.Event{
		EventName: event.EventName,
		Payload:   bs,
	}, nil
}
