package metadatautils

import (
	"context"
)

type Metadata map[string]string

func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, "metadata_key", md)
}

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value("metadata_key").(Metadata)
	if !ok {
		return nil, ok
	}

	newMetadata := map[string]string{}

	for k, v := range md {
		newMetadata[k] = v
	}

	return newMetadata, ok
}

func MergeContext(ctx context.Context, patch Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	md, _ := ctx.Value("metadata_key").(Metadata)

	cp := Metadata{}

	for k, v := range md {
		cp[k] = v
	}

	for k, v := range patch {
		_, ok := cp[k]
		if !ok || overwrite {
			cp[k] = v
		}
	}

	return context.WithValue(ctx, "metadata_key", cp)
}
