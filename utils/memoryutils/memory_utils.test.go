package memoryutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryBuffer(t *testing.T) {
	b := NewBuffer(BufferWithSize(10))

	b.Put("foo")

	entries := b.Get(1)

	val := entries[0].Value.(string)
	require.Equal(t, "foo", val)

	b = NewBuffer(BufferWithSize(10))

	for i := 0; i < 10; i++ {
		b.Put(i)
	}

	entries = b.Get(10)

	for i := 0; i < 10; i++ {
		val := entries[i].Value.(int)
		require.Equal(t, i, val)
	}

	for i := 0; i < 10; i++ {
		b.Put(i * 2)
	}

	entries = b.Get(10)

	for i := 0; i < 10; i++ {
		val := entries[i].Value.(int)
		require.Equal(t, i*2, val)
	}
}
