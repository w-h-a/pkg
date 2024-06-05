package log

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	size := 100
	
	service := "service.namespace"

	logger = NewLog(LogWithSize(size))
	require.Equal(t, size, logger.(*defaultLog).buffer.Options().Size)

	SetName(service)

	Info("foobar")

	SetLevel(LevelDebug)

	Debugf("foo %s", "bar")

	expected := []string{"foobar", "foo bar"}

	entries, err := logger.Read(ReadWithCount(len(expected)))
	require.NoError(t, err)

	for i, entry := range entries {
		message := entry.Message.(string)
		require.True(t, strings.Contains(message, expected[i]))
		require.True(t, strings.Contains(message, service))
	}
}
