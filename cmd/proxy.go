package cmd

import (
	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

func newProxyCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"p"},
		Short:   "Manage HTTP proxy server",
	}

	flags.AddCommand(
		newProxyStartCommand(),
		newProxyStopCommand(),
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

func newProxyStopCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:   "stop",
		Short: "Stop HTTP proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProxyStop(cmd)
		},
	}

	return flags
}

func runProxyStop(cmd *cobra.Command) error {
	return internal.Stop()
}
