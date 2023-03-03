package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/watcher"
)

type WatchArgs struct {
	mocksPath string

	domainsPath string
}

func watchCommand(subCommands ...*cobra.Command) *cobra.Command {
	var args WatchArgs

	command := &cobra.Command{
		Use: "watch",
		RunE: func(cmd *cobra.Command, _ []string) error {
			requests, err := createInputs(args)
			if err != nil {
				return err
			}

			runner := watcher.NewRunner(os.Stdout)

			if err = runner.Watch(cmd.Context(), requests...); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringVar(&args.mocksPath, "mocks", "", "Directory with mocks")
	command.Flags().StringVar(&args.domainsPath, "domains", "", "Directory with domain directories")

	command.AddCommand(subCommands...)

	return command
}

func createInputs(args WatchArgs) ([]watcher.WatchRequest, error) {
	if len(args.mocksPath) == 0 && len(args.domainsPath) == 0 {
		return nil, fmt.Errorf("watch command: at least one parameter must be set")
	}

	var watchers []watcher.WatchRequest

	if len(args.mocksPath) > 0 {
		watchers = append(watchers, watcher.WatchRequest{
			Path: args.mocksPath,
			Name: watcher.MocksWatcher,
		})
	}

	if len(args.domainsPath) > 0 {
		watchers = append(watchers, watcher.WatchRequest{
			Path: args.domainsPath,
			Name: watcher.DomainsWatcher,
		})
	}

	return watchers, nil
}
