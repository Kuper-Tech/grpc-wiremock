package schematomock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getStatusCodes(t *testing.T) {
	tests := []struct {
		name            string
		definedStatuses map[int]struct{}
		statusCodeRaw   string
		want            []int
	}{
		{definedStatuses: map[int]struct{}{}, statusCodeRaw: string(StatusCodeLevel5), want: statusCodesByLevel[StatusCodeLevel5]},
		{statusCodeRaw: string(StatusCodeLevel5), want: statusCodesByLevel[StatusCodeLevel5]},
		{definedStatuses: map[int]struct{}{101: {}, 103: {}}, statusCodeRaw: string(StatusCodeLevel1), want: []int{100, 102}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStatusCodes(tt.definedStatuses, tt.statusCodeRaw)
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
