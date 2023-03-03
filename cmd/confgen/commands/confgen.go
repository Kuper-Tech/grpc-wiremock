package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/utils/exec"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/confgen"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/runner"
)

type ConfgenArgs struct {
	wiremockPath string

	nginxPath       string
	supervisordPath string
}

func confgenCommand(subCommands ...*cobra.Command) *cobra.Command {
	var args ConfgenArgs

	command := &cobra.Command{
		Use: "confgen",
		RunE: func(cmd *cobra.Command, _ []string) error {
			commandRunner := runner.New(exec.New())

			generator := confgen.NewConfGenWithDefaultFs(
				args.wiremockPath, args.supervisordPath, commandRunner, os.Stdout)

			options, err := createOptions(args)
			if err != nil {
				return fmt.Errorf("create options: %w", err)
			}

			if err := generator.Generate(cmd.Context(), options...); err != nil {
				return fmt.Errorf("generate configs: %w", err)
			}

			return nil
		},
	}

	command.Flags().StringVar(&args.nginxPath, "nginx", environment.NginxConfigsPath, "Directory with NGINX config")
	command.Flags().StringVar(&args.wiremockPath, "wiremock-path", environment.DefaultWiremockConfigPath, "Directory with Wiremock config")
	command.Flags().StringVar(&args.supervisordPath, "supervisord", environment.SupervisordConfigsDirPath, "Directory with Supervisord config")

	command.AddCommand(subCommands...)

	return command
}

func createOptions(args ConfgenArgs) (confgen.Options, error) {
	if args.nginxPath == "" && args.supervisordPath == "" {
		return nil, fmt.Errorf("at least one argument must be set")
	}

	var options confgen.Options

	if args.nginxPath != "" {
		options = append(options,
			confgen.WithNGINX(args.nginxPath),
		)
	}

	if args.supervisordPath != "" {
		options = append(options,
			confgen.WithSupervisord(filepath.Join(args.supervisordPath, "mocks")),
		)
	}

	return options, nil
}
