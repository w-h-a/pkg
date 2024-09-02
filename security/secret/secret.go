package secret

type Secret interface {
	Options() SecretOptions
	GetSecret(name string, opts ...GetSecretOption) (map[string]string, error)
	String() string
}
