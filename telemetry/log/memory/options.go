package memory

import (
	"context"

	"github.com/w-h-a/pkg/telemetry/log"
)

type sizeKey struct{}

func LogWithSize(s int) log.LogOption {
	return func(o *log.LogOptions) {
		o.Context = context.WithValue(o.Context, sizeKey{}, s)
	}
}

func GetSizeFromContext(ctx context.Context) (int, bool) {
	s, ok := ctx.Value(sizeKey{}).(int)
	return s, ok
}
