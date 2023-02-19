package cmd

import (
	"fmt"
	"os"
	"path"
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
		newConfigRemoveCommand(),
	)

	return result
}

func newConfigListCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configurations",
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

func newConfigRemoveCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "remove <config>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigRemove(args[0])
		},
	}

	return result
}

func runConfigRemove(config string) error {
	configPath, err := internal.GetConfigPath()

	if err != nil {
		return err
	}

	count := 0

	for _, ext := range viper.SupportedExts {
		file := path.Join(configPath, config+"."+ext)
		info, err := os.Stat(file)

		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return err
			}
		}

		if !info.IsDir() {
			if err = os.Remove(file); err != nil {
				return err
			}
			count++
		}
	}

	if count == 0 {
		return fmt.Errorf("failed to find config %s", config)
	}

	return nil
}
