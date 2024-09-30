package http

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/serverv2"
	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/pkg/telemetry/log/memory"
)

func TestHttpServer(t *testing.T) {
	logger := memory.NewLog()

	log.SetLogger(logger)

	testResponse := "hello world"

	s := NewServer(
		serverv2.ServerWithAddress("127.0.0.1:0"),
	)

	err := s.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testResponse)
	}))
	require.NoError(t, err)

	server, ok := s.(*server)
	require.True(t, ok)

	err = server.Run()
	require.NoError(t, err)

	rsp, err := http.Get(fmt.Sprintf("http://%s/", s.Options().Address))
	require.NoError(t, err)

	defer rsp.Body.Close()

	bytes, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)
	require.Equal(t, testResponse, string(bytes))

	err = server.Stop()
	require.NoError(t, err)
}
