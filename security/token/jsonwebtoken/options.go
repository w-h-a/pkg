package jsonwebtoken

import (
	"context"

	"github.com/w-h-a/pkg/security/token"
)

type publicKeyKey struct{}
type privateKeyKey struct{}

func JwtWithPublicKey(key string) token.TokenOption {
	return func(o *token.TokenOptions) {
		o.Context = context.WithValue(o.Context, publicKeyKey{}, key)
	}
}

func GetPublicKeyFromContext(ctx context.Context) (string, bool) {
	k, ok := ctx.Value(publicKeyKey{}).(string)
	return k, ok
}

func JwtWithPrivateKey(key string) token.TokenOption {
	return func(o *token.TokenOptions) {
		o.Context = context.WithValue(o.Context, privateKeyKey{}, key)
	}
}

func GetPrivateKeyFromContext(ctx context.Context) (string, bool) {
	k, ok := ctx.Value(privateKeyKey{}).(string)
	return k, ok
}
