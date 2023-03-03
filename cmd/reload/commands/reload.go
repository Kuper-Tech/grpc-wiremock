package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	certgencmd "github.com/SberMarket-Tech/grpc-wiremock/cmd/certgen/commands"
	confgencmd "github.com/SberMarket-Tech/grpc-wiremock/cmd/confgen/commands"
)

func reloadCommand() *cobra.Command {
	command := &cobra.Command{
		Use: "reload",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			if err := certgencmd.CreateCommandRoot().ExecuteContext(ctx); err != nil {
				return fmt.Errorf("run certgen: %w", err)
			}

			if err := confgencmd.CreateCommandRoot().ExecuteContext(ctx); err != nil {
				return fmt.Errorf("run confgen: %w", err)
			}

			return nil
		},
	}

	return command
}
