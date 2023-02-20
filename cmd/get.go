package cmd

import (
	"os"

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
		Use:     "get <object>...",
		Aliases: []string{"g"},
		Args:    cobra.ArbitraryArgs,
		Short:   "Get objects from CyberArk",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args, params)
		},
	}

	result.Flags().StringVarP(&params.Config, "config", "c", "", "Config name")

	addConfigFlags(result.Flags(), &params)

	return result
}

func runGet(cmd *cobra.Command, args []string, params internal.Parameters) error {
	params.Objects = args

	if params.Config != "" {
		if _, err := loadConfig(cmd, params.Config); err != nil {
			return err
		}

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

	_, err := internal.GetPid()

	if os.IsNotExist(err) {
		if err = internal.Fork(); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return internal.NewClient(params).Run()
}
