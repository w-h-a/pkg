package memory

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/telemetry/log"
)

func TestLogger(t *testing.T) {
	size := 100

	service := "service.namespace"

	logger := NewLog(log.LogWithPrefix(service), LogWithSize(size))
	require.Equal(t, size, logger.(*defaultLog).buffer.Options().Size)

	log.Info("foobar")

	log.SetLevel(log.LevelDebug)

	log.Debugf("foo %s", "bar")

	expected := []string{"foobar", "foo bar"}

	entries, err := logger.Read(log.ReadWithCount(len(expected)))
	require.NoError(t, err)

	for i, entry := range entries {
		message := entry.Message.(string)
		require.True(t, strings.Contains(message, expected[i]))
		require.True(t, strings.Contains(message, service))
	}
}
