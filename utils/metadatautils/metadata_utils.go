package metadatautils

import (
	"context"
	"strings"
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
