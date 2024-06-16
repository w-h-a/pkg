package basictoken

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/security/token"
	"github.com/w-h-a/pkg/store/memory"
)

func TestGenerate(t *testing.T) {
	memStore := memory.NewStore()

	b := NewTokenProvider(
		BasicTokenWithStore(memStore),
	)

	_, err := b.Generate(
		token.GenerateWithId("test"),
	)
	require.NoError(t, err)

	records, err := memStore.List()
	require.NoError(t, err)
	require.Equal(t, 1, len(records))
}

func TestInspect(t *testing.T) {
	memStore := memory.NewStore()

	b := NewTokenProvider(
		BasicTokenWithStore(memStore),
	)

	t.Run("Valid token", func(t *testing.T) {
		md := map[string]string{"foo": "bar"}
		roles := []string{"admin"}
		id := "test"

		opts := []token.GenerateOption{
			token.GenerateWithId(id),
			token.GenerateWithRoles(roles...),
			token.GenerateWithMetadata(md),
		}

		token1, err := b.Generate(opts...)
		require.NoError(t, err)

		token2, err := b.Inspect(token1.AccessToken)
		require.NoError(t, err)

		require.Equal(t, id, token2.Id)
		require.True(t, len(token2.Roles) == len(roles))
		require.True(t, len(token2.Metadata) == len(md))
	})

	t.Run("Expired token", func(t *testing.T) {
		tk, err := b.Generate(
			token.GenerateWithId("foo"),
			token.GenerateWithExpiry(-10*time.Second),
		)
		require.NoError(t, err)

		_, err = b.Inspect(tk.AccessToken)
		require.True(t, err == token.ErrInvalidToken)
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := b.Inspect("invalid")
		require.True(t, err == token.ErrInvalidToken)
	})
}
