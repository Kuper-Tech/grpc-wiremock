package certificates

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configopener"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

var (
	osfs = afero.NewOsFs()

	staticFS = static.FromEmbed()

	projectDir = filepath.Join(fsutils.CurrentDir(), "../../..")
)

func Test_certsGen_Generate(t *testing.T) {
	tests := []struct {
		name string

		fs afero.Fs

		domains       []string
		commonDomains []string
		wiremockPath  string
		want          []string
	}{
		{
			fs: osfs,

			domains:       []string{"com", "ru"},
			commonDomains: []string{"localhost", "grpc-wiremock", "google.com"},
			want:          []string{"localhost", "grpc-wiremock", "google.com", "awesome.ru", "awesome.com", "push-sender.ru", "push-sender.com"},

			wiremockPath: filepath.Join(projectDir, "static/tests/data/supervisord/two-services"),
		},
		{
			fs: staticFS,

			domains:       []string{"com", "ru"},
			commonDomains: []string{"localhost", "grpc-wiremock", "google.com"},
			want:          []string{"localhost", "grpc-wiremock", "google.com", "awesome.ru", "awesome.com"},

			wiremockPath: filepath.Join(projectDir, "static/tests/data/supervisord/one-service"),
		},
		{
			fs: staticFS,

			domains:       []string{"com", "ru"},
			commonDomains: []string{"localhost", "grpc-wiremock"},
			want:          []string{"localhost", "grpc-wiremock"},

			wiremockPath: filepath.Join(projectDir, "static/tests/data/supervisord/empty"),
		},
		{
			fs: staticFS,

			domains:       []string{"com", "ru"},
			commonDomains: []string{"localhost", "grpc-wiremock", "google.com"},
			want:          []string{"localhost", "grpc-wiremock", "google.com"},

			wiremockPath: filepath.Join(projectDir, "static/tests/data/supervisord/empty-dir"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &certsGen{opener: configopener.New(tt.fs, tt.wiremockPath), fs: tt.fs}
			got, err := g.collectDomains(tt.commonDomains, tt.domains)
			require.NoError(t, err)
			require.ElementsMatch(t, got, tt.want)
		})
	}
}
