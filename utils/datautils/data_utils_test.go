package datautils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringify(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		data := "TestString"
		actual, err := Stringify(data)
		require.NoError(t, err)
		require.Equal(t, []byte(data), actual)
	})

	t.Run("byte array type", func(t *testing.T) {
		data := []byte("TestString")
		actual, err := Stringify(data)
		require.NoError(t, err)
		require.Equal(t, data, actual)
	})

	t.Run("map[string]string type", func(t *testing.T) {
		data := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		result, err := Stringify(data)
		require.NoError(t, err)

		var actual map[string]string
		err = json.Unmarshal(result, &actual)
		require.NoError(t, err)

		require.EqualValues(t, data, actual)
	})
}
