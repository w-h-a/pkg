package errorutils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			Id:     "test",
			Code:   500,
			Detail: "Internal server error",
			Status: http.StatusText(500),
		},
	}

	for _, e := range testData {
		err := NewError(e.Id, e.Detail, e.Code)
		require.Equal(t, e.Error(), err.Error())

		parsedError := ParseError(err.Error())
		require.NotNil(t, parsedError)
		require.Equal(t, e.Id, parsedError.Id)
		require.Equal(t, e.Detail, parsedError.Detail)
		require.Equal(t, e.Code, parsedError.Code)
		require.Equal(t, e.Status, parsedError.Status)
	}
}
