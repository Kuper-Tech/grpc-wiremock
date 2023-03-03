package testutils

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

var fs = static.FromEmbed()

func TestReadTestsStatuses(t *testing.T) {
	tests := []struct {
		name string
		fs   afero.Fs
		path string
		want TestToCases
	}{
		{
			fs:   fs,
			path: "tests-statuses.yml",
			want: map[string]testCasesMaps{
				"generation": {
					Skip:    map[string]struct{}{"simple-1-flaky": {}},
					Fail:    map[string]struct{}{"simple-2-fail": {}},
					Success: map[string]struct{}{"simple-success": {}},
				},
				"run-some-fancy-test-name": {
					Skip:    map[string]struct{}{"flaky-test-1": {}},
					Fail:    map[string]struct{}{},
					Success: map[string]struct{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadTestsStatuses(tt.fs, tt.path)
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}
