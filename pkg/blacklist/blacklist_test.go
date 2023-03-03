package blacklist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDeliveredWithProtoc(t *testing.T) {
	tests := []struct {
		name string

		value string

		want bool
	}{
		{value: "google.golang.org/protobuf/reflect/protoreflect", want: true},
		{value: "google.golang.org/protobuf/runtime/protoimpl", want: true},
		{value: "google/protobuf/wrapper.proto", want: true},
		{value: "google.golang.org/grpc/codes", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDeliveredWithProtoc(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}
