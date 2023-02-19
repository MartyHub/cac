package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func newConfigCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
	}

	result.AddCommand(
		newConfigListCommand(),
	)

	return result
}

func newConfigListCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "list",
		Short: "List configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigList(cmd)
		},
	}

	return result
}

func runConfigList(cmd *cobra.Command) error {
	configPath, err := internal.GetConfigPath()

	if err != nil {
		return err
	}

	entries, err := os.ReadDir(configPath)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())

			if slices.Contains(viper.SupportedExts, ext[1:]) {
				if _, err = cmd.OutOrStdout().Write(
					[]byte(fmt.Sprintln(strings.TrimSuffix(entry.Name(), ext))),
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
