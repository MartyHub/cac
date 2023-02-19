package cmd

import (
	"time"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	certFileName       = "cert-file"
	keyFileName        = "key-file"
	hostName           = "host"
	appIdName          = "app-id"
	safeName           = "safe"
	jsonName           = "json"
	maxConnectionsName = "max-connections"
	maxTriesName       = "max-tries"
	timeoutName        = "timeout"
	waitName           = "wait"
)

func newGetCommand() *cobra.Command {
	params := internal.NewParameters()
	result := &cobra.Command{
		Use:   "get <object>...",
		Short: "Get objects from CyberArk",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args, params)
		},
	}

	result.Flags().StringVarP(&params.Config, "config", "c", "", "Config name")

	result.Flags().StringVar(&params.CertFile, certFileName, "", "Certificate file")
	result.Flags().StringVar(&params.KeyFile, keyFileName, "", "Key file")

	result.Flags().StringVar(&params.Host, hostName, "", "CyberArk CCP REST Web Service Host")
	result.Flags().StringVar(&params.AppId, appIdName, "", "CyberArk Application Id")
	result.Flags().StringVar(&params.Safe, safeName, "", "CyberArk Safe")

	result.Flags().BoolVar(&params.Json, jsonName, false, "JSON output")
	result.Flags().IntVar(&params.MaxConns, maxConnectionsName, 4, "Max connections")
	result.Flags().IntVar(&params.MaxTries, maxTriesName, 3, "Max tries")
	result.Flags().DurationVar(&params.Timeout, timeoutName, 30*time.Second, "Timeout")
	result.Flags().DurationVar(&params.Wait, waitName, 100*time.Millisecond, "Wait before retry")

	cobra.CheckErr(viper.BindPFlag(certFileName, result.Flags().Lookup(certFileName)))
	cobra.CheckErr(viper.BindPFlag(keyFileName, result.Flags().Lookup(keyFileName)))

	cobra.CheckErr(viper.BindPFlag(hostName, result.Flags().Lookup(hostName)))
	cobra.CheckErr(viper.BindPFlag(appIdName, result.Flags().Lookup(appIdName)))
	cobra.CheckErr(viper.BindPFlag(safeName, result.Flags().Lookup(safeName)))

	cobra.CheckErr(viper.BindPFlag(jsonName, result.Flags().Lookup(jsonName)))
	cobra.CheckErr(viper.BindPFlag(maxConnectionsName, result.Flags().Lookup(maxConnectionsName)))
	cobra.CheckErr(viper.BindPFlag(maxTriesName, result.Flags().Lookup(maxTriesName)))
	cobra.CheckErr(viper.BindPFlag(timeoutName, result.Flags().Lookup(timeoutName)))
	cobra.CheckErr(viper.BindPFlag(waitName, result.Flags().Lookup(waitName)))

	return result
}

func runGet(cmd *cobra.Command, args []string, params internal.Parameters) error {
	params.Objects = args

	if params.Config != "" {
		loadConfig(params.Config)

		params.CertFile = viper.GetString(certFileName)
		params.KeyFile = viper.GetString(keyFileName)

		params.AppId = viper.GetString(appIdName)
		params.Host = viper.GetString(hostName)
		params.Safe = viper.GetString(safeName)

		params.Json = viper.GetBool(jsonName)
		params.MaxConns = viper.GetInt(maxConnectionsName)
		params.MaxTries = viper.GetInt(maxTriesName)
		params.Timeout = viper.GetDuration(timeoutName)
		params.Wait = viper.GetDuration(waitName)
	}

	if err := params.Validate(); err != nil {
		return err
	}

	return internal.NewClient(params).Run()
}

func loadConfig(config string) {
	configPath, err := internal.GetConfigPath()

	cobra.CheckErr(err)

	viper.AddConfigPath(configPath)
	viper.SetConfigName(config)
	viper.SetConfigPermissions(0o600)

	cobra.CheckErr(viper.ReadInConfig())
	cobra.CheckErr(internal.CheckConfigFilePermissions(viper.ConfigFileUsed()))

}
