package env

import (
	"os"

	"github.com/w-h-a/pkg/security/secret"
)

type envSecret struct {
	options secret.SecretOptions
}

func (s *envSecret) Options() secret.SecretOptions {
	return s.options
}

func (s *envSecret) GetSecret(name string, opts ...secret.GetSecretOption) (map[string]string, error) {
	options := secret.NewGetSecretOptions(opts...)

	key := options.Prefix + name

	value := os.Getenv(key)

	return map[string]string{
		name: value,
	}, nil
}

func (s *envSecret) String() string {
	return "env"
}

func NewSecret(opts ...secret.SecretOption) secret.Secret {
	options := secret.NewSecretOptions(opts...)

	s := &envSecret{
		options: options,
	}

	return s
}
