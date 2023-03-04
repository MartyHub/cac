package cmd

import (
	"strings"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	params := internal.NewParameters()
	result := &cobra.Command{
		Use:     "get <object>...",
		Aliases: []string{"g"},
		Args:    cobra.ArbitraryArgs,
		Short:   "Get accounts from CyberArk",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd, args, params)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if cache, err := internal.NewCache(params.Config); err == nil {
				return cache.SortedObjects(strings.ToLower(toComplete)), cobra.ShellCompDirectiveNoFileComp
			}

			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}

	result.Flags().StringVarP(&params.Config, "config", "c", "", "Config name")
	_ = result.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		result, err := getConfigs(toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return result, cobra.ShellCompDirectiveNoFileComp
	})
	result.Flags().StringVarP(&params.Output, "output", "o", "", "Generate files in given output path")

	addConfigFlags(result, &params)

	return result
}

func runGet(cmd *cobra.Command, args []string, params internal.Parameters) error {
	params, err := applyConfig(cmd, params)
	if err != nil {
		return err
	}

	params.Objects = args

	if err := params.Validate(); err != nil {
		return err
	}

	client, err := internal.NewClient(params)

	if err != nil {
		return err
	}

	return client.Run()
}
