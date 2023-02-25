package cmd

import (
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
	expiryName         = "expiry"
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

		if err := viper.Unmarshal(&params); err != nil {
			return err
		}
	}

	if err := params.Validate(); err != nil {
		return err
	}

	client, err := internal.NewClient(params)

	if err != nil {
		return err
	}

	return client.Run()
}
