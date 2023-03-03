package confgen

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

const (
	nginxOption       = "nginx"
	supervisordOption = "supervisord"
)

var staticFS = static.FromEmbed()

type option struct {
	name       string
	enable     bool
	outputPath string

	prepareFunctions []preparator
}

type preparator func(afero.Fs, string) error

func withOption(name, outputPath string) option {
	return option{enable: true, name: name, outputPath: outputPath}
}

func withOptionAndPreparation(name, outputPath string, prepareF ...preparator) option {
	return option{enable: true, name: name, outputPath: outputPath, prepareFunctions: prepareF}
}

func WithNGINX(outputPath string) option {
	return withOption(nginxOption, outputPath)
}

func WithSupervisord(outputPath string, prepareF ...preparator) option {
	return withOptionAndPreparation(supervisordOption, outputPath, prepareF...)
}

func WithDefaultSupervisordConfig(fs afero.Fs, path string) error {
	const templatePath = "supervisord/default-config-path"

	if err := fsutils.CopyDir(staticFS, fs, templatePath, path, false); err != nil {
		return fmt.Errorf("prepare default supervisord config: %w", err)
	}

	return nil
}

type Options []option

func (o Options) Paths() []string {
	var paths []string

	for _, option := range o {
		paths = append(paths, option.outputPath)
	}

	return paths
}
