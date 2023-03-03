package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPorts_Allocate(t *testing.T) {
	tests := []struct {
		name string
		p    Ports
		want int
	}{
		{p: Ports{8000, 8002, 8003}, want: 8004},
		{p: Ports{8007}, want: 8008},
		{p: Ports{}, want: 8000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.Allocate()
			require.Equal(t, tt.want, got)
		})
	}
}
