package kubernetes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/runtime"
)

type testcase struct {
	Token         string
	Method        string
	URL           string
	Header        map[string]string
	NewReqWrapper func(options *runtime.RuntimeOptions) *request
}

var tests = []testcase{
	{
		NewReqWrapper: func(options *runtime.RuntimeOptions) *request {
			return newRequest(options).get().setResource("service")
		},
		Token:  "my fake token",
		Method: "GET",
		URL:    "/api/v1/namespaces/default/services/",
	},
	{
		NewReqWrapper: func(options *runtime.RuntimeOptions) *request {
			return newRequest(options).get().setResource("pod").setParams(&params{labelSelector: map[string]string{"foo": "bar"}})
		},
		Token:  "my fake token",
		Method: "GET",
		URL:    "/api/v1/namespaces/default/pods/?labelSelector=foo%3Dbar",
	},
	{
		NewReqWrapper: func(options *runtime.RuntimeOptions) *request {
			return newRequest(options).get().setResource("deployment").setParams(&params{labelSelector: map[string]string{"foo": "bar"}})
		},
		Token:  "my fake token",
		Method: "GET",
		URL:    "/apis/apps/v1/namespaces/default/deployments/?labelSelector=foo%3Dbar",
	},
}

var wrappedHandler = func(test *testcase, t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer "+test.Token, auth)
		require.Equal(t, test.Method, r.Method)
		require.Equal(t, test.URL, r.URL.RequestURI())
		w.WriteHeader(http.StatusOK)
	})
}

func TestKubernetes(t *testing.T) {
	for _, test := range tests {
		ts := httptest.NewServer(wrappedHandler(&test, t))
		req := test.NewReqWrapper(&runtime.RuntimeOptions{
			Host:        ts.URL,
			Namespace:   "default",
			BearerToken: test.Token,
			Client:      &http.Client{},
		})
		res := req.do()
		require.NoError(t, res.getError())
		ts.Close()
	}
}
