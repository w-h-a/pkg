package retryutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	results := []time.Duration{
		0 * time.Second,
		100 * time.Millisecond,
		600 * time.Millisecond,
		1900 * time.Millisecond,
		4300 * time.Millisecond,
		7900 * time.Millisecond,
	}

	for i := 0; i < 5; i++ {
		d := ExponentialBackoff(i)
		require.Equal(t, d, results[i])
	}
}
