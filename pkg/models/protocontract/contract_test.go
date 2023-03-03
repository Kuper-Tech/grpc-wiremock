package protocontract

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

func TestSourceFileTypeFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want types.SourceFileType
	}{
		{path: "/some/test/path/file.proto", want: types.ProtoType},
		{path: "/some/test/path/file.yaml", want: types.OpenAPIType},
		{path: "/some/test/path/file.yml", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.SourceFileTypeFromPath(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSourceFileType_Is(t *testing.T) {
	tests := []struct {
		name string
		s    types.SourceFileType
		that string
		want bool
	}{
		{s: types.OpenAPIType, that: "yaml", want: true},
		{s: types.ProtoType, that: "proto", want: true},
		{s: types.OpenAPIType, that: "yml", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Is(tt.that)
			assert.Equal(t, tt.want, got)
		})
	}
}
