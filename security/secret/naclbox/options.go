package naclbox

import (
	"context"

	"github.com/w-h-a/pkg/security/secret"
)

type publicKeyKey struct{}
type privateKeyKey struct{}

func NaclboxWithKeys(public, private []byte) secret.SecretOption {
	return func(o *secret.SecretOptions) {
		o.Context = context.WithValue(o.Context, publicKeyKey{}, public)
		o.Context = context.WithValue(o.Context, privateKeyKey{}, private)
	}
}

func EncryptWithPublicKey(public []byte) secret.EncryptOption {
	return func(o *secret.EncryptOptions) {
		o.Context = context.WithValue(o.Context, publicKeyKey{}, public)
	}
}

func DecryptWithPublicKey(public []byte) secret.DecryptOption {
	return func(o *secret.DecryptOptions) {
		o.Context = context.WithValue(o.Context, publicKeyKey{}, public)
	}
}

func GetPublicKeyFromContext(ctx context.Context) ([]byte, bool) {
	b, ok := ctx.Value(publicKeyKey{}).([]byte)
	return b, ok
}

func GetPrivateKeyFromContext(ctx context.Context) ([]byte, bool) {
	b, ok := ctx.Value(privateKeyKey{}).([]byte)
	return b, ok
}
