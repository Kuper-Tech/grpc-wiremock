package commands

import "github.com/spf13/cobra"

func CreateCommandRoot() *cobra.Command {
	return watchCommand()
}
