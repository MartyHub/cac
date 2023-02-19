package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "cac",
		Short: "Simple CyberArk Central Credentials Provider REST client",
	}

	result.AddCommand(
		newConfigCommand(),
		newProxyCommand(),
		newGetCommand(),
		newVersionCommand(),
	)

	return result
}
