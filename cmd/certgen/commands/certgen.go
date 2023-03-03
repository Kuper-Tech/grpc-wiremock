package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/certgen"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
)

type CertgenArgs struct {
	supervisordPath string

	output string
}

func certgenCommand(subCommands ...*cobra.Command) *cobra.Command {
	var args CertgenArgs

	command := &cobra.Command{
		Use: "certgen",
		RunE: func(cmd *cobra.Command, _ []string) error {
			generator := certgen.NewCertsGenWithDefaultFs(args.output, args.supervisordPath, os.Stdout)

			if err := generator.Generate(cmd.Context()); err != nil {
				return fmt.Errorf("generate certs: %w", err)
			}

			return nil
		},
	}

	command.Flags().StringVar(&args.supervisordPath, "supervisord-path", environment.SupervisordConfigsDirPath, "Directory with Supervisord config")
	command.Flags().StringVar(&args.output, "output", environment.DefaultCertificatesPath, "Directory with certificates")

	command.AddCommand(subCommands...)

	return command
}
