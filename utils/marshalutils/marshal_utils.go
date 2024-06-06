package marshalutils

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	DefaultMarshalers = map[string]Marshaler{
		"application/grpc":         protoMarshaler{},
		"application/grpc+json":    jsonMarshaler{},
		"application/grpc+proto":   protoMarshaler{},
		"application/json":         jsonMarshaler{},
		"application/octet-stream": protoMarshaler{},
		"application/proto":        protoMarshaler{},
		"application/protobuf":     protoMarshaler{},
	}
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Name() string
}

type jsonMarshaler struct{}

func (jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonMarshaler) Unmarshal(d []byte, v interface{}) error {
	if len(d) == 0 {
		return nil
	}

	protoMessage, ok := v.(proto.Message)
	if ok {
		return protojson.Unmarshal(d, protoMessage)
	}

	return json.Unmarshal(d, v)
}

func (jsonMarshaler) Name() string {
	return "json"
}

type protoMarshaler struct{}

func (protoMarshaler) Marshal(v interface{}) ([]byte, error) {
	protoMessage, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal: %v is not a proto message", v)
	}

	return proto.Marshal(protoMessage)
}

func (protoMarshaler) Unmarshal(d []byte, v interface{}) error {
	if len(d) == 0 {
		return nil
	}

	protoMessage, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal: %v is not a proto message", v)
	}

	return proto.Unmarshal(d, protoMessage)
}

func (protoMarshaler) Name() string {
	return "proto"
}
