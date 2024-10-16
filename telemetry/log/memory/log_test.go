package memory

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func TestLogger(t *testing.T) {
	size := 100

	service := "service.namespace"

	logger := NewLog(log.LogWithPrefix(service), log.LogWithLevel(log.LevelDebug), LogWithBuffer(memoryutils.NewBuffer(memoryutils.BufferWithSize(size))))
	require.Equal(t, size, logger.(*memoryLog).buffer.Options().Size)

	log.SetLogger(logger)

	log.Info("foobar")

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
