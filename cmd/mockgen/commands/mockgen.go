package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/mockgen"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

type MockgenArgs struct {
	domain       string
	inputPath    string
	wiremockPath string
	contractType string
}

func mockgenCommand(subCommands ...*cobra.Command) *cobra.Command {
	var args MockgenArgs

	command := &cobra.Command{
		Use: "mockgen",
		RunE: func(cmd *cobra.Command, _ []string) error {
			contractType, err := createContractType(args.contractType)
			if err != nil {
				return err
			}

			gen := mockgen.NewMocksGenWithDefaultFs(args.inputPath, args.wiremockPath, os.Stdout)
			return gen.Generate(cmd.Context(), contractType, args.domain)
		},
	}

	command.Flags().StringVar(&args.inputPath, "input", "", "Directory with contract files")
	command.Flags().StringVar(&args.wiremockPath, "wiremock-path", environment.DefaultWiremockConfigPath, "Directory with Wiremock config")
	command.Flags().StringVar(&args.domain, "domain", "awesome-service", "Name of the domain")
	command.Flags().StringVar(&args.contractType, "type", "", "Select `openapi` or `proto` file type")

	if err := command.MarkFlagRequired("input"); err != nil {
		log.Fatalf("mark flag 'input' required: %s", err)
	}

	if err := command.MarkFlagRequired("type"); err != nil {
		log.Fatalf("mark flag 'type' required: %s", err)
	}

	command.AddCommand(subCommands...)

	return command
}

func createContractType(value string) (types.SourceFileType, error) {
	switch value {
	case "proto":
		return types.ProtoType, nil
	case "openapi":
		return types.OpenAPIType, nil
	}

	return "", fmt.Errorf("unsupported contract type")
}
