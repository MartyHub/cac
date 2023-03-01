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
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

const (
	aliasesName        = "aliases"
	appIdName          = "app-id"
	certFileName       = "cert-file"
	expiryName         = "expiry"
	hostName           = "host"
	jsonName           = "json"
	keyFileName        = "key-file"
	maxConnectionsName = "max-connections"
	maxTriesName       = "max-tries"
	outputName         = "output"
	safeName           = "safe"
	skipVerifyName     = "skip-verify"
	timeoutName        = "timeout"
	waitName           = "wait"
)

func newConfigCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
	}

	result.AddCommand(
		newConfigListCommand(),
		newConfigRemoveCommand(),
		newConfigSetCommand(),
	)

	return result
}

func newConfigListCommand() *cobra.Command {
	verbose := false
	result := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigList(cmd, verbose)
		},
	}

	result.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return result
}

func runConfigList(cmd *cobra.Command, verbose bool) error {
	configInfos, err := getConfigInfos()
	if err != nil {
		return err
	}

	for _, configInfo := range configInfos {
		cmd.Println(configInfo.name)

		if verbose {
			if err = printConfig(cmd, configInfo.file); err != nil {
				return err
			}
		}
	}

	return nil
}

func printConfig(cmd *cobra.Command, file string) error {
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	for _, name := range []string{
		aliasesName,
		appIdName,
		certFileName,
		expiryName,
		hostName,
		jsonName,
		keyFileName,
		maxConnectionsName,
		maxTriesName,
		safeName,
		skipVerifyName,
		timeoutName,
		waitName,
	} {
		cmd.Println("  ", fmt.Sprintf("%-15s", name), "=", viper.Get(name))
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
		ValidArgsFunction: completeConfig,
	}

	return result
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
			return runConfigSet(cmd, args[0], params.Aliases)
		},
		ValidArgsFunction: completeConfig,
	}

	addConfigFlags(result, &params)

	result.Flags().StringSliceVar(&params.Aliases, aliasesName, []string{}, "Aliases")
	_ = result.RegisterFlagCompletionFunc(aliasesName, cobra.NoFileCompletions)

	result.Flags().DurationVar(&params.Expiry, expiryName, 12*time.Hour, "Cache expiry")

	return result
}

func runConfigSet(cmd *cobra.Command, config string, aliases []string) error {
	configPath, err := loadConfig(cmd, config)

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigFile(path.Join(configPath, config+".yaml"))
		} else {
			return err
		}
	}

	return viper.WriteConfig()
}

func applyConfig(cmd *cobra.Command, params internal.Parameters) (internal.Parameters, error) {
	if params.Config == "" {
		return params, nil
	}

	configHome, err := loadConfig(cmd, params.Config)
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = loadConfigAlias(configHome, params.Config)
	}
	if err != nil {
		return params, err
	}

	err = viper.Unmarshal(&params)

	return params, err
}

func loadConfigAlias(configHome string, alias string) error {
	entries, err := os.ReadDir(configHome)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())

			if slices.Contains(viper.SupportedExts, ext[1:]) {
				viper.SetConfigFile(filepath.Join(configHome, entry.Name()))

				if err = viper.ReadInConfig(); err == nil {
					if slices.Contains(viper.GetStringSlice(aliasesName), alias) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("failed to find config with alias %s", alias)
}

func loadConfig(cmd *cobra.Command, config string) (string, error) {
	configHome, err := internal.GetConfigHome()

	if err != nil {
		return "", err
	}

	viper.AddConfigPath(configHome)
	viper.SetConfigName(config)
	viper.SetConfigPermissions(0o600)

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		return configHome, err
	}

	err = viper.ReadInConfig()

	return configHome, err
}

func addConfigFlags(cmd *cobra.Command, params *internal.Parameters) {
	cmd.Flags().StringVar(&params.CertFile, certFileName, "", "Certificate file")
	_ = cmd.MarkFlagFilename(certFileName, "cer", "cert", "crt", "pem")

	cmd.Flags().StringVar(&params.KeyFile, keyFileName, "", "Key file")
	_ = cmd.MarkFlagFilename(certFileName, "cer", "cert", "crt", "key", "pem")

	cmd.Flags().StringVar(&params.Host, hostName, "", "CyberArk CCP REST Web Service Host")
	_ = cmd.RegisterFlagCompletionFunc(hostName, cobra.NoFileCompletions)

	cmd.Flags().StringVar(&params.AppId, appIdName, "", "CyberArk Application Id")
	_ = cmd.RegisterFlagCompletionFunc(appIdName, cobra.NoFileCompletions)

	cmd.Flags().StringVar(&params.Safe, safeName, "", "CyberArk Safe")
	_ = cmd.RegisterFlagCompletionFunc(safeName, cobra.NoFileCompletions)

	cmd.Flags().BoolVar(&params.Json, jsonName, false, "JSON output")
	cmd.Flags().BoolVar(&params.SkipVerify, skipVerifyName, false, "Skip server certificate verification")
	cmd.Flags().IntVar(&params.MaxConns, maxConnectionsName, 4, "Max connections")
	cmd.Flags().IntVar(&params.MaxTries, maxTriesName, 3, "Max tries")
	cmd.Flags().DurationVar(&params.Timeout, timeoutName, 30*time.Second, "Timeout")
	cmd.Flags().DurationVar(&params.Wait, waitName, 100*time.Millisecond, "Wait before retry")
}

func completeConfig(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	result, err := getConfigs(toComplete)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

type configInfo struct {
	name, file string
}

func getConfigInfos() ([]configInfo, error) {
	configHome, err := internal.GetConfigHome()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(configHome)
	if err != nil {
		return nil, err
	}

	var result []configInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())

			if slices.Contains(viper.SupportedExts, ext[1:]) {
				config := strings.TrimSuffix(entry.Name(), ext)

				result = append(result, configInfo{
					name: config,
					file: filepath.Join(configHome, entry.Name()),
				})
			}
		}
	}

	return result, nil
}

func getConfigs(prefix string) ([]string, error) {
	configInfos, err := getConfigInfos()
	if err != nil {
		return nil, err
	}

	var result []string

	for _, configInfo := range configInfos {
		if strings.HasPrefix(configInfo.name, prefix) {
			result = append(result, configInfo.name)
		}
	}

	return result, nil
}
