package ssm

import (
	"context"

	"github.com/w-h-a/pkg/security/secret"
)

type ssmClientKey struct{}

func SsmWithSsmClient(c SsmClient) secret.SecretOption {
	return func(o *secret.SecretOptions) {
		o.Context = context.WithValue(o.Context, ssmClientKey{}, c)
	}
}

func GetSsmClientFromContext(ctx context.Context) (SsmClient, bool) {
	c, ok := ctx.Value(ssmClientKey{}).(SsmClient)
	return c, ok
}
