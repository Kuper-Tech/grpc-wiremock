package strutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToPackageName(t *testing.T) {
	tests := []struct {
		name string

		input string

		want string
	}{
		{input: "github.com/foo/bar/baugi", want: "baugi"},
		{input: "github.com/foo/bar/baugiWithSmthg", want: "baugi_with_smthg"},
		{input: "github.com/foo/bar/baugi.with.point", want: "baugi_with_point"},
		{input: "baugi.with.point", want: "baugi_with_point"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToPackageName(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string

		input string

		want string
	}{
		{input: "baugi", want: "baugi"},
		{input: "baugiWithSmthg", want: "baugi_with_smthg"},
		{input: "baugi.with.point", want: "baugi_with_point"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToPackageName(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUniqueAndSorted(t *testing.T) {
	tests := []struct {
		name string

		values []string

		want []string
	}{
		{values: []string{"g", "a", "A", "a", "b"}, want: []string{"A", "a", "b", "g"}},
		{values: []string{"foo", "foo"}, want: []string{"foo"}},
		{values: []string{"foo", "bar"}, want: []string{"bar", "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqueAndSorted(tt.values...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string

		input string

		want string
	}{
		{
			input: "some_string",
			want:  "SomeString",
		},
		{
			input: "",
			want:  "",
		},
		{
			input: "SomeString",
			want:  "SomeString",
		},
		{
			input: "Some String",
			want:  "SomeString",
		},
		{
			input: "Some-String",
			want:  "SomeString",
		},
		{
			input: "Some-string",
			want:  "SomeString",
		},
		{
			input: "some-string",
			want:  "SomeString",
		},
		{
			input: "some/string",
			want:  "SomeString",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
