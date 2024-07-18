package secret

type Secret interface {
	Options() SecretOptions
	Encrypt(bytes []byte, opts ...EncryptOption) ([]byte, error)
	Decrypt(bytes []byte, opts ...DecryptOption) ([]byte, error)
	String() string
}
