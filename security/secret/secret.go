package secret

type Secret interface {
	Options() SecretOptions
	Encrypt(messageToSend []byte, opts ...EncryptOption) ([]byte, error)
	Decrypt(receivedMessage []byte, opts ...DecryptOption) ([]byte, error)
	String() string
}
