package naclbox

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/w-h-a/pkg/security/secret"
	"github.com/w-h-a/pkg/telemetry/log"
	"golang.org/x/crypto/nacl/box"
)

const (
	keyLength = 32
)

type naclbox struct {
	options    secret.SecretOptions
	publicKey  [keyLength]byte
	privateKey [keyLength]byte
}

func (s *naclbox) Options() secret.SecretOptions {
	return s.options
}

func (s *naclbox) Encrypt(messageToSend []byte, opts ...secret.EncryptOption) ([]byte, error) {
	options := secret.NewEncryptOptions(opts...)

	public, ok := GetPublicKeyFromContext(options.Context)
	if !ok || len(public) != keyLength {
		return []byte{}, fmt.Errorf("a length %d public key of the recipient of this encrypted message must be provided", keyLength)
	}

	var recipientPublicKey [keyLength]byte

	copy(recipientPublicKey[:], public)

	var nonce [24]byte

	if _, err := rand.Reader.Read(nonce[:]); err != nil {
		return []byte{}, fmt.Errorf("failed to obtain random nonce from crypto/rand: %v", err)
	}

	return box.Seal(nonce[:], messageToSend, &nonce, &recipientPublicKey, &s.privateKey), nil
}

func (s *naclbox) Decrypt(receivedMessage []byte, opts ...secret.DecryptOption) ([]byte, error) {
	options := secret.NewDecryptOptions(opts...)

	public, ok := GetPublicKeyFromContext(options.Context)
	if !ok || len(public) != keyLength {
		return []byte{}, fmt.Errorf("a length %d public key of the sender of this encrypted message must be provided", keyLength)
	}

	var senderPublicKey [keyLength]byte

	copy(senderPublicKey[:], public)

	var nonce [24]byte

	copy(nonce[:], receivedMessage[:24])

	msg, ok := box.Open(nil, receivedMessage[24:], &nonce, &senderPublicKey, &s.privateKey)
	if !ok {
		return []byte{}, errors.New("failed to decrypt the message")
	}

	return msg, nil
}

func (s *naclbox) String() string {
	return "naclbox"
}

func NewSecret(opts ...secret.SecretOption) secret.Secret {
	options := secret.NewSecretOptions(opts...)

	s := &naclbox{
		options: options,
	}

	public, ok1 := GetPublicKeyFromContext(options.Context)

	private, ok2 := GetPrivateKeyFromContext(options.Context)

	if !ok1 || !ok2 || len(public) != keyLength || len(private) != keyLength {
		log.Fatalf("a public and private key of length %d must both be provided", keyLength)
	}

	copy(s.publicKey[:], public)
	copy(s.privateKey[:], private)

	return s
}
