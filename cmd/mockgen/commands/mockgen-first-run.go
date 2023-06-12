package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	mockgen_first_run "github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/mockgen-first-run"
)

type MockgenFirstRunArgs struct {
	DomainsPath  string
	WiremockPath string
}

func mockgenFirstRun(subCommands ...*cobra.Command) *cobra.Command {
	var args MockgenFirstRunArgs

	command := &cobra.Command{
		Use: "first-run",
		RunE: func(cmd *cobra.Command, _ []string) error {
			gen := mockgen_first_run.NewMocksGenWithDefaultFs(args.DomainsPath, args.WiremockPath, os.Stdout)
			return gen.GenerateForEachDomain(cmd.Context())
		},
	}

	command.Flags().StringVar(&args.DomainsPath, "domains-path", "", "Directory with domain directories")
	command.Flags().StringVar(&args.WiremockPath, "wiremock-path", "", "Directory with Wiremock config")

	if err := command.MarkFlagRequired("domains-path"); err != nil {
		log.Fatalf("mark flag 'domains-path' required: %s", err)
	}

	if err := command.MarkFlagRequired("wiremock-path"); err != nil {
		log.Fatalf("mark flag 'wiremock-path' required: %s", err)
	}

	command.AddCommand(subCommands...)

	return command
}
