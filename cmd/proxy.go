package cmd

import (
	"fmt"
	"os"

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
		newProxyStatusCommand(),
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

func newProxyStatusCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:   "status",
		Short: "Display HTTP proxy server status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProxyStatus(cmd)
		},
	}

	return flags
}

func runProxyStatus(cmd *cobra.Command) error {
	pid, err := internal.GetPid()

	if os.IsNotExist(err) {
		_, err := cmd.OutOrStdout().Write([]byte("Proxy is not running\n"))

		return err
	}

	if err != nil {
		return err
	}

	_, err = cmd.OutOrStdout().Write([]byte(fmt.Sprintf("Proxy is running with PID %d\n", pid)))

	return err
}
