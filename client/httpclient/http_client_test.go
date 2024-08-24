package httpclient

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/utils/marshalutils"
)

type Payload struct {
	Seq  int64  `json:"seq,omitempty"`
	Data string `json:"data,omitempty"`
}

func TestHttpClient(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:3000")
	require.NoError(t, err)
	defer l.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "expect post method", 500)
			return
		}

		ct := r.Header.Get("content-type")
		marshaler, ok := marshalutils.DefaultMarshalers[ct]
		if !ok {
			http.Error(w, "marshaler not found", 500)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		payload := &Payload{}

		if err := marshaler.Unmarshal(b, payload); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err = marshaler.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	go http.Serve(l, mux)

	c := NewClient()

	for i := 0; i < 10; i++ {
		payload := &Payload{
			Seq:  int64(i),
			Data: fmt.Sprintf("message %d", i),
		}

		req := c.NewRequest(
			client.RequestWithNamespace("test"),
			client.RequestWithName("test"),
			client.RequestWithPort(3000),
			client.RequestWithMethod("/foo/bar"),
			client.RequestWithUnmarshaledRequest(payload),
		)

		rsp := &Payload{}

		err := c.Call(context.Background(), req, rsp, client.CallWithAddress("127.0.0.1:3000"))

		require.NoError(t, err)

		require.Equal(t, payload.Seq, rsp.Seq)
	}
}
