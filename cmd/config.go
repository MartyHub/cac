package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

const (
	aliasesName        = "aliases"
	appIDName          = "app-id"
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

	extJSON = ".json"
)

const rw = 0o600

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
		RunE: func(cmd *cobra.Command, _ []string) error {
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

func readConfigFile(file string) (internal.Config, error) {
	var result internal.Config

	data, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)

	return result, err
}

func readConfig(name string) (internal.Config, error) {
	var result internal.Config

	configHome, err := internal.GetConfigHome()
	if err != nil {
		return result, err
	}

	result, err = readConfigFile(filepath.Join(configHome, name+extJSON))
	if err == nil {
		return result, nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	return readConfigAlias(configHome, name)
}

func readConfigAlias(configHome, alias string) (internal.Config, error) {
	var result internal.Config

	entries, err := os.ReadDir(configHome)
	if err != nil {
		return result, err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), extJSON) {
			continue
		}

		file := filepath.Join(configHome, entry.Name())

		result, err = readConfigFile(file)
		if err != nil {
			return result, err
		}

		if internal.Contains(result.Aliases, alias) {
			return result, nil
		}
	}

	return result, internal.NewError(nil, "failed to find config %q", alias)
}

func writeConfig(file string, cfg internal.Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, rw)
}

func printConfig(cmd *cobra.Command, file string) error {
	cfg, err := readConfigFile(file)
	if err != nil {
		return err
	}

	cmd.Println(cfg.String())

	return nil
}

func newConfigRemoveCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "remove <config>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a configuration",
		RunE: func(_ *cobra.Command, args []string) error {
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

	return os.Remove(filepath.Join(configHome, config+extJSON))
}

func newConfigSetCommand() *cobra.Command {
	cfg := internal.NewConfig()
	result := &cobra.Command{
		Use:   "set <config>",
		Args:  cobra.ExactArgs(1),
		Short: "Add or update a configuration",
		RunE: func(_ *cobra.Command, args []string) error {
			return runConfigSet(args[0], cfg)
		},
		ValidArgsFunction: completeConfig,
	}

	result.Flags().StringSliceVar(&cfg.Aliases, aliasesName, []string{}, "Aliases")
	_ = result.RegisterFlagCompletionFunc(aliasesName, cobra.NoFileCompletions)

	result.Flags().StringVar(&cfg.AppID, appIDName, "", "CyberArk Application Id")
	_ = result.RegisterFlagCompletionFunc(appIDName, cobra.NoFileCompletions)

	result.Flags().StringVar(&cfg.CertFile, certFileName, "", "Certificate file")
	_ = result.MarkFlagFilename(certFileName, "cer", "cert", "crt", "pem")

	result.Flags().DurationVar(&cfg.Expiry, expiryName, cfg.Expiry, "Cache expiry")

	result.Flags().StringVar(&cfg.Host, hostName, "", "CyberArk CCP REST Web Service Host")
	_ = result.RegisterFlagCompletionFunc(hostName, cobra.NoFileCompletions)

	result.Flags().StringVar(&cfg.KeyFile, keyFileName, "", "Key file")
	_ = result.MarkFlagFilename(keyFileName, "cer", "cert", "crt", "key", "pem")

	result.Flags().IntVar(&cfg.MaxConns, maxConnectionsName, cfg.MaxConns, "Max connections")
	result.Flags().IntVar(&cfg.MaxTries, maxTriesName, cfg.MaxTries, "Max tries")

	result.Flags().StringVar(&cfg.Safe, safeName, "", "CyberArk Safe")
	_ = result.RegisterFlagCompletionFunc(safeName, cobra.NoFileCompletions)

	result.Flags().BoolVar(&cfg.SkipVerify, skipVerifyName, false, "Skip server certificate verification")
	result.Flags().DurationVar(&cfg.Timeout, timeoutName, cfg.Timeout, "Timeout")
	result.Flags().DurationVar(&cfg.Wait, waitName, cfg.Wait, "Wait before retry")

	return result
}

func runConfigSet(name string, cfg internal.Config) error {
	configHome, err := internal.GetConfigHome()
	if err != nil {
		return err
	}

	file := filepath.Join(configHome, name+".json")

	existCfg, err := readConfigFile(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return writeConfig(file, cfg)
		}

		return err
	}

	return writeConfig(file, existCfg.Overwrite(cfg))
}

func completeConfig(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		if strings.HasSuffix(entry.Name(), ".json") {
			result = append(result, configInfo{
				name: strings.TrimSuffix(entry.Name(), ".json"),
				file: filepath.Join(configHome, entry.Name()),
			})
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
