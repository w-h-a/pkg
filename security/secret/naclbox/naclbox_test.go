package naclbox

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
)

func TestNaclbox(t *testing.T) {
	alicePublicKey, alicePrivateKey, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	bobPublicKey, bobPrivateKey, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)

	alice := NewSecret(NaclboxWithKeys(alicePublicKey[:], alicePrivateKey[:]))
	require.Equal(t, "naclbox", alice.String())

	bob := NewSecret(NaclboxWithKeys(bobPublicKey[:], bobPrivateKey[:]))
	require.Equal(t, "naclbox", bob.String())

	aliceSecretMsg := []byte("23 is number 1")

	_, err = alice.Encrypt(aliceSecretMsg)
	require.Error(t, err)

	bobPub, ok := GetPublicKeyFromContext(bob.Options().Context)
	require.True(t, ok)

	encrypted, err := alice.Encrypt(aliceSecretMsg, EncryptWithPublicKey(bobPub))
	require.NoError(t, err)

	_, err = bob.Decrypt(encrypted)
	require.Error(t, err)

	alicePub, ok := GetPublicKeyFromContext(alice.Options().Context)
	require.True(t, ok)

	decrypted, err := bob.Decrypt(encrypted, DecryptWithPublicKey(alicePub))
	require.NoError(t, err)

	require.Equal(t, aliceSecretMsg, decrypted)
}
