package metadatautils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	md := Metadata{
		"Foo": "bar",
	}

	ctx := NewContext(context.TODO(), md)

	meta, ok := FromContext(ctx)
	require.True(t, ok)
	require.Equal(t, md["Foo"], meta["Foo"])
	require.Equal(t, 1, len(meta))
}

func TestMergeContext(t *testing.T) {
	type args struct {
		existing  Metadata
		append    Metadata
		overwrite bool
	}

	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "matching key, overwrite false",
			args: args{
				existing:  Metadata{"Foo": "bar", "Baz": "test1"},
				append:    Metadata{"Baz": "test2"},
				overwrite: false,
			},
			want: Metadata{"Foo": "bar", "Baz": "test1"},
		},
		{
			name: "matching key, overwrite true",
			args: args{
				existing:  Metadata{"Foo": "bar", "Baz": "test1"},
				append:    Metadata{"Baz": "test2"},
				overwrite: true,
			},
			want: Metadata{"Foo": "bar", "Baz": "test2"},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, _ := FromContext(MergeContext(NewContext(context.TODO(), c.args.existing), c.args.append, c.args.overwrite))
			require.Equal(t, c.want, got)
		})
	}
}
