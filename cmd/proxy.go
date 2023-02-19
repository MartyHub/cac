package cmd

import (
	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

func newProxyCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"d"},
		Short:   "Manage HTTP proxy server",
	}

	flags.AddCommand(
		newProxyStartCommand(),
	)

	return flags
}

func newProxyStartCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:    "start",
		Hidden: true,
		Short:  "Start HTTP proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProxyStart(cmd)
		},
	}

	return flags
}

func runProxyStart(cmd *cobra.Command) error {
	return internal.Start()
}
