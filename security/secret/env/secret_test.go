package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/security/secret"
)

func TestSecret(t *testing.T) {
	s := NewSecret()

	t.Setenv("TEST_SECRET", "secret1")
	require.Equal(t, "secret1", os.Getenv("TEST_SECRET"))

	t.Run("Get", func(t *testing.T) {
		rsp, err := s.GetSecret("TEST_SECRET")
		require.NoError(t, err)
		require.NotNil(t, rsp)
		require.Equal(t, "secret1", rsp["TEST_SECRET"])
	})
}

func TestSecretWithPrefix(t *testing.T) {
	s := NewSecret()

	t.Setenv("TEST_SECRET", "test1")
	t.Setenv("test_secret", "test2")
	t.Setenv("FOOP_SECRET", "test3")
	require.Equal(t, "test1", os.Getenv("TEST_SECRET"))
	require.Equal(t, "test2", os.Getenv("test_secret"))
	require.Equal(t, "test3", os.Getenv("FOOP_SECRET"))

	t.Run("Get", func(t *testing.T) {
		rsp, err := s.GetSecret("SECRET", secret.GetSecretWithPrefix("TEST_"))
		require.NoError(t, err)
		require.NotNil(t, rsp)
		require.Len(t, rsp, 1)
		require.Equal(t, "test1", rsp["SECRET"])
	})
}
