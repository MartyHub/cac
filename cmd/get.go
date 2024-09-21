package cmd

import (
	"strings"

	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	params := internal.NewParameters()
	result := &cobra.Command{
		Use:     "get <config> <account>...",
		Aliases: []string{"g"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Get accounts from CyberArk",
		RunE: func(_ *cobra.Command, args []string) error {
			return runGet(args, params)
		},
		ValidArgsFunction: func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return completeConfig(cmd, args, toComplete)
			}

			return completeAccount(args[0], args[1:], toComplete)
		},
	}

	result.Flags().BoolVarP(&params.JSON, jsonName, "j", false, "Output JSON")
	result.Flags().StringVarP(&params.Output, outputName, "o", "", "Generate files in given output path")

	return result
}

func completeAccount(config string, exclusions []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cache, err := internal.NewDBCache()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	defer cache.Close()

	accounts, err := cache.SortedAccounts(config, strings.ToLower(toComplete), exclusions)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	result := make([]string, len(accounts))

	for i, acct := range accounts {
		result[i] = acct.Object
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func runGet(args []string, params internal.Parameters) error {
	var err error

	params.CfgName = args[0]
	params.Objects = args[1:]

	params.Config, err = readConfig(params.CfgName)
	if err != nil {
		return err
	}

	if err = params.Validate(); err != nil {
		return err
	}

	client, err := internal.NewClient(params)
	if err != nil {
		return err
	}

	return client.Run()
}
