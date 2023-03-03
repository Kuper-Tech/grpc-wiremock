package configs

import (
	"context"
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

type Configuer interface {
	GenerateConfig(values Values) error
}

type Reloader interface {
	ReloadConfig(ctx context.Context) error
}

type Values struct {
	Domain string
	Port   string
	Root   string
}

type runner struct {
	fs     afero.Fs
	config config.Wiremock
}

func NewRunner(fs afero.Fs, config config.Wiremock) *runner {
	return &runner{fs: fs, config: config}
}

func (g *runner) RunConfiguers(configuers ...Configuer) error {
	for _, service := range g.config.Services {
		values := Values{
			Domain: service.Name,
			Root:   service.RootDir,
			Port:   fmt.Sprint(service.Port),
		}

		for _, configuer := range configuers {
			if err := configuer.GenerateConfig(values); err != nil {
				return fmt.Errorf("generate config: %w", err)
			}
		}
	}

	return nil
}

func (g *runner) RunReloaders(ctx context.Context, reloaders ...Reloader) error {
	for _, reloader := range reloaders {
		if err := reloader.ReloadConfig(ctx); err != nil {
			return fmt.Errorf("reload config: %w", err)
		}
	}

	return nil
}
