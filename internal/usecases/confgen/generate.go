package confgen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/configs"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/configs/nginx"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/configs/supervisord"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/renderer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configopener"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configsync"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

//go:generate mockery --exported --name=commandRunner --filename=generate_mock.go
type commandRunner interface {
	Run(ctx context.Context, cmd string, args ...string) error
}

type confGen struct {
	WiremockPath string

	SupervisordPath string

	Logger io.Writer

	afero.Fs
	commandRunner
}

var configuerToTemplatePath = map[string]string{
	nginxOption:       "proxy-nginx/files",
	supervisordOption: "supervisord/files",
}

var reloaderToCommand = map[string][]string{
	nginxOption:       {"sudo", "nginx", "-s", "reload"},
	supervisordOption: {"supervisord", "ctl", "reload", "-c", environment.SupervisordMainConfigPath},
}

func NewConfGenWithDefaultFs(wiremockPath, supervisordPath string, commandRunner commandRunner, logger io.Writer) confGen {
	return confGen{
		WiremockPath:    wiremockPath,
		SupervisordPath: supervisordPath, Logger: logger,
		commandRunner: commandRunner, Fs: afero.NewOsFs(),
	}
}

func (g *confGen) Generate(ctx context.Context, options ...option) error {
	if len(options) == 0 {
		return fmt.Errorf("at least one option should be set")
	}

	wiremockConfig, err := g.syncWiremockConfig()
	if err != nil {
		return fmt.Errorf("create wiremock config: %w", err)
	}

	if err := g.preparePaths(options...); err != nil {
		return fmt.Errorf("prepare paths: %w", err)
	}

	configuers, err := g.createConfiguers(static.FromEmbed(), options...)
	if err != nil {
		return fmt.Errorf("create configures: %w", err)
	}

	reloaders, err := g.createReloaders(options...)
	if err != nil {
		return fmt.Errorf("create configures: %w", err)
	}

	runner := configs.NewRunner(g.Fs, wiremockConfig)

	if err = runner.RunConfiguers(configuers...); err != nil {
		return fmt.Errorf("configure: %w", err)
	}

	if err = runner.RunReloaders(ctx, reloaders...); err != nil {
		return fmt.Errorf("reload: %w", err)
	}

	return nil
}

func (g *confGen) syncWiremockConfig() (config.Wiremock, error) {
	wiremockConfig, err := configopener.New(g.Fs, g.SupervisordPath).Open()
	if err != nil {
		log.Printf("open wiremock config: %s", err)
		return config.Wiremock{}, handleOpenerErrors(err)
	}

	targetWiremockConfig, err := configsync.SyncWiremockConfig(g.Fs, wiremockConfig, g.WiremockPath)
	if err != nil {
		var e syscall.Errno
		if !errors.As(err, &e) {
			return config.Wiremock{}, fmt.Errorf("sync wiremock config: %w", err)
		}
	}

	return targetWiremockConfig, nil
}

func (g *confGen) applyOptions(options ...option) error {
	for _, option := range options {
		if option.prepareFunctions == nil {
			continue
		}

		for _, prepF := range option.prepareFunctions {
			if err := prepF(g.Fs, option.outputPath); err != nil {
				return fmt.Errorf("prepare func: %w", err)
			}
		}
	}

	return nil
}

func (g *confGen) createConfiguer(staticFS afero.Fs, option option) (configs.Configuer, error) {
	templatePath, exists := configuerToTemplatePath[option.name]
	if !exists {
		return nil, fmt.Errorf("template path doesn't exist: %s", option.name)
	}

	configRenderer, err := renderer.New(staticFS, templatePath)
	if err != nil {
		return nil, fmt.Errorf("create renderer: %w", err)
	}

	switch option.name {
	case nginxOption:
		return nginx.Configuer{
			Fs:         g.Fs,
			Renderer:   configRenderer,
			OutputPath: option.outputPath,
		}, nil
	case supervisordOption:
		return supervisord.Configuer{
			Fs:         g.Fs,
			Renderer:   configRenderer,
			OutputPath: option.outputPath,
		}, nil
	}

	return nil, fmt.Errorf("unknown option name: %s", option.name)
}

func (g *confGen) createReloader(option option) (configs.Reloader, error) {
	command, exists := reloaderToCommand[option.name]
	if !exists {
		return nil, fmt.Errorf("command doesn't exist: %s", option.name)
	}

	switch option.name {
	case nginxOption:
		return nginx.Reloader{
			Command: command,
			Runner:  g.commandRunner,
		}, nil
	case supervisordOption:
		return supervisord.Reloader{
			Command: command,
			Runner:  g.commandRunner,
		}, nil
	}

	return nil, fmt.Errorf("unknown option name: %s", option.name)
}

func (g *confGen) createConfiguers(staticFS afero.Fs, options ...option) ([]configs.Configuer, error) {
	var configuers []configs.Configuer

	for _, opt := range options {
		configuer, err := g.createConfiguer(staticFS, opt)
		if err != nil {
			return nil, fmt.Errorf("create configuer: %w", err)
		}

		configuers = append(configuers, configuer)
	}

	return configuers, nil
}

func (g *confGen) createReloaders(options ...option) ([]configs.Reloader, error) {
	var reloaders []configs.Reloader

	for _, opt := range options {
		reloader, err := g.createReloader(opt)
		if err != nil {
			return nil, fmt.Errorf("create reloader: %w", err)
		}

		reloaders = append(reloaders, reloader)
	}

	return reloaders, nil
}

func (g *confGen) preparePaths(options ...option) error {
	configPaths := Options(options).Paths()

	if err := g.cleanConfigs(configPaths); err != nil {
		return fmt.Errorf("clean configs: %w", err)
	}

	if err := g.applyOptions(options...); err != nil {
		return fmt.Errorf("apply options: %w", err)
	}

	return nil
}

func (g *confGen) cleanConfigs(paths []string) error {
	for _, configPath := range paths {
		if err := environment.CleanConfigs(g.Fs, configPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return fmt.Errorf("clean configs: %w", err)
		}
	}

	return nil
}

func handleOpenerErrors(err error) error {
	if errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, config.EmptyWiremockConfigErr) {
		return nil
	}

	return fmt.Errorf("generate: %w", err)
}
