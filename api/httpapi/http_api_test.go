package httpapi

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/api"
)

func TestHttpApi(t *testing.T) {
	testResponse := "hello world"

	a := NewApi(
		api.ApiWithAddress("127.0.0.1:0"),
	)

	a.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testResponse)
	}))

	server, ok := a.(*httpApi)
	require.True(t, ok)

	err := server.start()
	require.NoError(t, err)

	rsp, err := http.Get(fmt.Sprintf("http://%s/", a.Options().Address))
	require.NoError(t, err)

	defer rsp.Body.Close()

	bytes, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)
	require.Equal(t, testResponse, string(bytes))

	err = server.stop()
	require.NoError(t, err)
}
