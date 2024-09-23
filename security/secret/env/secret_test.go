package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
