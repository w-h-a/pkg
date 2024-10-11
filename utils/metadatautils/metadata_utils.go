package metadatautils

import (
	"context"
	"net/http"
	"strings"

	"github.com/w-h-a/pkg/telemetry/log"
)

type metadataKey struct{}

type Metadata map[string]string

func NewContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	cp := Metadata{}

	for k, v := range md {
		cp[strings.ToLower(k)] = v
	}

	return context.WithValue(ctx, metadataKey{}, cp)
}

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
	}

	newMetadata := map[string]string{}

	for k, v := range md {
		newMetadata[strings.ToLower(k)] = v
	}

	return newMetadata, ok
}

func SetContext(ctx context.Context, k, v string) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = Metadata{}
	}

	md[strings.ToLower(k)] = v

	return context.WithValue(ctx, metadataKey{}, md)
}

func GetContext(ctx context.Context, k string) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", ok
	}

	val, ok := md[strings.ToLower(k)]
	log.Infof("GET CONTEXT KEY %s, VALUE %s, md %#+v", k, val, md)

	return val, ok
}

func MergeContext(ctx context.Context, patch Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := ctx.Value(metadataKey{}).(Metadata)

	cp := Metadata{}

	for k, v := range md {
		cp[strings.ToLower(k)] = v
	}

	for k, v := range patch {
		_, ok := cp[strings.ToLower(k)]
		if !ok || overwrite {
			cp[strings.ToLower(k)] = v
		}
	}

	return context.WithValue(ctx, metadataKey{}, cp)
}

func RequestToContext(r *http.Request) context.Context {
	ctx := context.Background()

	md := Metadata{}

	for k, v := range r.Header {
		log.Infof("KEY %s, VALUE %s", k, v)
		md[strings.ToLower(k)] = strings.Join(v, ",")
	}

	return NewContext(ctx, md)
}
