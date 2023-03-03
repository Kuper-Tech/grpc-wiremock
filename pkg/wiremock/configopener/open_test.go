package configopener

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

var osfs = afero.NewOsFs()

var projectPath = filepath.Join(fsutils.CurrentDir(), "../../..")

func Test_opener_Open(t *testing.T) {
	tests := []struct {
		name    string
		fs      afero.Fs
		path    string
		want    config.Wiremock
		wantErr bool
	}{
		{
			fs:   osfs,
			path: filepath.Join(projectPath, "static/tests/data/supervisord/simple"),
			want: config.Wiremock{},
		},
		{
			fs:   osfs,
			path: filepath.Join(projectPath, "static/tests/data/supervisord/empty-dir"),
			want: config.Wiremock{},
		},
		{
			fs:   osfs,
			path: filepath.Join(projectPath, "static/tests/data/supervisord/with-includes"),
			want: config.Wiremock{Services: []config.Service{{Port: 8000, Name: "awesome", RootDir: "/home/mock/awesome"}}},
		},
		{
			fs:   osfs,
			path: filepath.Join(projectPath, "static/tests/data/supervisord/two-services"),
			want: config.Wiremock{Services: []config.Service{
				{Port: 8000, Name: "awesome", RootDir: "/home/mock/awesome"},
				{Port: 8001, Name: "push-sender", RootDir: "/home/mock/push-sender"},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &opener{fs: tt.fs, path: tt.path}
			got, err := o.Open()
			require.NoError(t, err)
			require.ElementsMatch(t, got.Services, tt.want.Services)
		})
	}
}

func Test_convertEnvs(t *testing.T) {
	tests := []struct {
		name string
		envs []string
		want map[string]string
	}{
		{envs: []string{}, want: map[string]string{}},
		{envs: []string{"A:c"}, want: map[string]string{}},
		{envs: []string{"A=c=d"}, want: map[string]string{}},
		{envs: []string{"A=c"}, want: map[string]string{"A": "c"}},
		{envs: []string{"A=c", "ac=123"}, want: map[string]string{"A": "c", "ac": "123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertEnvs(tt.envs)
			require.Equal(t, tt.want, got)
		})
	}
}
