package grpc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToHandlerMethod(t *testing.T) {
	type testCase struct {
		input   string
		handler string
		method  string
		err     bool
	}

	grpcFormattedMethods := []testCase{
		{"/Foo/Bar", "Foo", "Bar", false},
		{"a.package.Foo/Bar", "", "", true},
	}

	for _, test := range grpcFormattedMethods {
		handler, method, err := ToHandlerMethod(test.input)
		if err != nil && test.err == true {
			continue
		}

		if err != nil && test.err == false {
			t.Fatalf("unexpected err %v for %+v", err, test)
		}

		if test.err == true && err == nil {
			t.Fatalf("expected error for %+v but got %s.%s", test, handler, method)
		}

		require.Equal(t, test.handler, handler)

		require.Equal(t, test.method, method)
	}
}
