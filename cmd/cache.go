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
	caches, err := internal.GetCaches("")
	if err != nil {
		return err
	}

	for _, cache := range caches {
		cmd.Println(cache)

		if verbose {
			if err = printCache(cmd, cache); err != nil {
				return err
			}
		}
	}

	return nil
}

func printCache(cmd *cobra.Command, config string) error {
	cache, err := internal.NewCache(config)
	if err != nil {
		return err
	}

	for _, object := range cache.SortedObjects() {
		cmd.Println("  ", object, "=", cache.Accounts[object].Value)
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

			result, err := internal.GetCaches(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return result, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return result
}

func runCacheRemove(config string) error {
	cache, err := internal.NewCache(config)
	if err != nil {
		return err
	}

	return cache.Remove()
}
