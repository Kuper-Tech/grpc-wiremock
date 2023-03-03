package sliceutils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstOf(t *testing.T) {
	tests := []struct {
		name string

		slice []string

		want string
	}{
		{slice: []string{"foo", "bar"}, want: "foo"},
		{slice: []string{}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FirstOf(tt.slice)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceToMap(t *testing.T) {
	type testCase[T comparable] struct {
		name string

		slice []T

		want map[T]struct{}
	}
	tests := []testCase[string]{
		{slice: []string{"1", "1", "2"}, want: map[string]struct{}{"1": {}, "2": {}}},
		{slice: []string{""}, want: map[string]struct{}{"": {}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToMap(tt.slice)
			assert.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
