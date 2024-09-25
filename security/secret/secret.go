package secret

type Secret interface {
	Options() SecretOptions
	GetSecret(key string) (map[string]string, error)
	String() string
}
