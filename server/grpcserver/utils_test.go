package grpcserver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToControllerHandler(t *testing.T) {
	type testCase struct {
		input      string
		controller string
		handler    string
		err        bool
	}

	grpcFormattedMethods := []testCase{
		{"/Foo/Bar", "Foo", "Bar", false},
		{"a.package.Foo/Bar", "", "", true},
	}

	for _, test := range grpcFormattedMethods {
		controller, handler, err := ToControllerHandler(test.input)
		if err != nil && test.err == true {
			continue
		}

		if err != nil && test.err == false {
			t.Fatalf("unexpected err %v for %+v", err, test)
		}

		if test.err == true && err == nil {
			t.Fatalf("expected error for %+v but got %s.%s", test, controller, handler)
		}

		require.Equal(t, test.controller, controller)

		require.Equal(t, test.handler, handler)
	}
}
