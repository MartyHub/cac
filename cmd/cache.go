package cmd

import (
	"os"
	"path/filepath"
	"strings"

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
	stateHome, err := internal.GetStateHome()

	if err != nil {
		return err
	}

	entries, err := os.ReadDir(stateHome)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := filepath.Ext(entry.Name())

			if ext == ".json" {
				config := strings.TrimSuffix(entry.Name(), ext)
				cmd.Println(config)

				if verbose {
					if err = printCache(cmd, config); err != nil {
						return err
					}
				}
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
