package watcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getDomainByPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{path: "/home/mock/awesome/mappings/mock.json", want: "awesome"},
		{path: "/home/mock/awesome-domain/mappings/mock.json", want: "awesome-domain"},
		{path: "/home/mock/awesome-domain/__files/mock.json", want: "awesome-domain"},
		{path: "/Users/test-user1/Desktop/awesome-project/deps/services/dependency/mappings/mock.json", want: "dependency"},

		{path: "awesome/__files/file.ext", want: "awesome"},
		{path: "/awesome/__files", want: "awesome"},
		{path: "q_123-test-not-so(awesome/mappings", want: "q_123-test-not-so(awesome"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDomainByPath(tt.path)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
