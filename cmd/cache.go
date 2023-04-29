package cmd

import (
	"github.com/MartyHub/cac/internal"
	"github.com/spf13/cobra"
)

func newCacheCommand() *cobra.Command {
	result := &cobra.Command{
		Use:   "cache",
		Short: "Manage caches",
	}

	result.AddCommand(
		newCacheListCommand(),
		newCacheRemoveCommand(),
	)

	return result
}

func newCacheListCommand() *cobra.Command {
	verbose := false
	result := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List caches",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCacheList(cmd, verbose)
		},
	}

	result.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return result
}

func runCacheList(cmd *cobra.Command, verbose bool) error {
	cache, err := internal.NewDBCache()
	if err != nil {
		return err
	}

	defer cache.Close()

	configs, err := cache.Configs("")
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		cmd.Println(cfg)

		if verbose {
			if err = printCache(cmd, cache, cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

func printCache(cmd *cobra.Command, cache internal.DBCache, config string) error {
	accounts, err := cache.SortedAccounts(config, "", nil)
	if err != nil {
		return err
	}

	for _, acct := range accounts {
		cmd.Println("  ", acct.Object, "=", acct.Value)
	}

	return nil
}

func newCacheRemoveCommand() *cobra.Command {
	result := &cobra.Command{
		Use:     "remove <cache>",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		Short:   "Remove a cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCacheRemove(args[0])
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			cache, err := internal.NewDBCache()
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			defer cache.Close()

			result, err := cache.Configs(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return result, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return result
}

func runCacheRemove(config string) error {
	cache, err := internal.NewDBCache()
	if err != nil {
		return err
	}

	defer cache.Close()

	return cache.RemoveAll(config)
}
