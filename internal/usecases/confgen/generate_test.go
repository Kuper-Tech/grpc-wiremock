package confgen

import (
	"context"
	"io"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/confgen/mocks"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

var (
	osfs        = afero.NewOsFs()
	emptyCtx    = mock.AnythingOfType("*context.emptyCtx")
	projectPath = filepath.Join(fsutils.CurrentDir(), "../../..")
	outputPath  = filepath.Join(projectPath, "tests")
)

func Test_confGen_Generate(t *testing.T) {
	tests := []struct {
		WiremockPath string
		Logger       io.Writer
		Fs           afero.Fs
		name         string
	}{
		{Fs: osfs, WiremockPath: filepath.Join(projectPath, "static/tests/data/supervisord/empty-dir"), name: "empty-dir"},
		{Fs: osfs, WiremockPath: filepath.Join(projectPath, "static/tests/data/supervisord/simple"), name: "simple"},
		{Fs: osfs, WiremockPath: filepath.Join(projectPath, "static/tests/data/supervisord/with-includes"), name: "with-includes"},
		{Fs: osfs, WiremockPath: filepath.Join(projectPath, "static/tests/data/supervisord/two-services"), name: "two-services"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := mocks.NewCommandRunner(t)
			runner.On("Run", emptyCtx, "supervisord", "ctl", "reload", "-c", environment.SupervisordMainConfigPath).Return(nil)

			outputPath := filepath.Join(outputPath, tt.name)
			g := &confGen{WiremockPath: tt.WiremockPath, Logger: tt.Logger, Fs: tt.Fs, commandRunner: runner}

			err := g.Generate(context.Background(), WithSupervisord(outputPath, WithDefaultSupervisordConfig))
			require.NoError(t, err)
		})
	}
}
