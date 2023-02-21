package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

func newConfigCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c"},
		Short:   "Manage configurations",
	}

	flags.AddCommand(
		newConfigListCommand(),
		newConfigRemoveCommand(),
		newConfigSetCommand(),
	)

	return flags
}

func newConfigListCommand() *cobra.Command {
	flags := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigList(cmd)
		},
	}

	return flags
}

func runConfigList(cmd *cobra.Command) error {
	configHome, err := internal.GetConfigHome()

	if err != nil {
		return err
	}

	entries, err := os.ReadDir(configHome)

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
	flags := &cobra.Command{
		Use:     "remove <config>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigRemove(args[0])
		},
	}

	return flags
}

func runConfigRemove(config string) error {
	configHome, err := internal.GetConfigHome()

	if err != nil {
		return err
	}

	count := 0

	for _, ext := range viper.SupportedExts {
		file := path.Join(configHome, config+"."+ext)
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

func newConfigSetCommand() *cobra.Command {
	params := internal.NewParameters()
	result := &cobra.Command{
		Use:   "set <config>",
		Args:  cobra.ExactArgs(1),
		Short: "Add or update a configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			params.Config = args[0]

			return runConfigSet(cmd, params)
		},
	}

	addConfigFlags(result.Flags(), &params)

	return result
}

func runConfigSet(cmd *cobra.Command, params internal.Parameters) error {
	configPath, err := loadConfig(cmd, params.Config)

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.WriteConfigAs(path.Join(configPath, params.Config+".properties"))
		} else {
			return err
		}
	}

	return viper.WriteConfig()
}

func loadConfig(cmd *cobra.Command, config string) (string, error) {
	configHome, err := internal.GetConfigHome()

	if err != nil {
		return "", err
	}

	viper.AddConfigPath(configHome)
	viper.SetConfigName(config)
	viper.SetConfigPermissions(0o600)

	// First read config
	err = viper.ReadInConfig()

	if err != nil {
		return configHome, err
	}

	// Then bind flags
	err = bindConfigFlags(cmd)

	return configHome, err
}

func addConfigFlags(flags *pflag.FlagSet, params *internal.Parameters) {
	flags.StringVar(&params.CertFile, certFileName, "", "Certificate file")
	flags.StringVar(&params.KeyFile, keyFileName, "", "Key file")

	flags.StringVar(&params.Host, hostName, "", "CyberArk CCP REST Web Service Host")
	flags.StringVar(&params.AppId, appIdName, "", "CyberArk Application Id")
	flags.StringVar(&params.Safe, safeName, "", "CyberArk Safe")

	flags.BoolVar(&params.Json, jsonName, false, "JSON output")
	flags.IntVar(&params.MaxConns, maxConnectionsName, 4, "Max connections")
	flags.IntVar(&params.MaxTries, maxTriesName, 3, "Max tries")
	flags.DurationVar(&params.Timeout, timeoutName, 30*time.Second, "Timeout")
	flags.DurationVar(&params.Wait, waitName, 100*time.Millisecond, "Wait before retry")
}

func bindConfigFlags(cmd *cobra.Command) error {
	for _, name := range []string{
		certFileName,
		keyFileName,
		hostName,
		appIdName,
		safeName,
		jsonName,
		maxConnectionsName,
		maxTriesName,
		timeoutName,
		waitName,
	} {
		if err := viper.BindPFlag(name, cmd.Flags().Lookup(name)); err != nil {
			return err
		}
	}

	return nil
}
