package secret

type Secret interface {
	Options() SecretOptions
	GetSecret(name string) (map[string]string, error)
	String() string
}
