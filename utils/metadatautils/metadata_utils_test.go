package metadatautils

import (
	"context"
	"net/http"
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
	require.Equal(t, md["Foo"], meta["foo"])
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
			want: Metadata{"foo": "bar", "baz": "test1"},
		},
		{
			name: "matching key, overwrite true",
			args: args{
				existing:  Metadata{"Foo": "bar", "Baz": "test1"},
				append:    Metadata{"Baz": "test2"},
				overwrite: true,
			},
			want: Metadata{"foo": "bar", "baz": "test2"},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, _ := FromContext(MergeContext(NewContext(context.TODO(), c.args.existing), c.args.append, c.args.overwrite))
			require.Equal(t, c.want, got)
		})
	}
}

func TestRequestToContext(t *testing.T) {
	testData := []struct {
		request *http.Request
		expect  Metadata
	}{
		{
			&http.Request{
				Header: http.Header{
					"Foo1": []string{"bar"},
					"Foo2": []string{"bar", "baz"},
				},
			},
			Metadata{
				"foo1": "bar",
				"foo2": "bar,baz",
			},
		},
	}

	for _, d := range testData {
		ctx := RequestToContext(d.request)
		md, ok := FromContext(ctx)
		require.True(t, ok)

		for k, v := range d.expect {
			val := md[k]
			require.Equal(t, v, val)
		}
	}
}
