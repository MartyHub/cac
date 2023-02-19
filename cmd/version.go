package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Value is injected during the build
var Version = "unknown"

func newVersionCommand() *cobra.Command {
	full := false
	result := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(cmd, full)
		},
	}

	result.Flags().BoolVarP(&full, "full", "f", false, "also display Git revision and modification time")

	return result
}

func runVersion(cmd *cobra.Command, full bool) error {
	var err error

	if full {
		_, err = cmd.OutOrStdout().Write([]byte(getFullVersion()))
	} else {
		_, err = cmd.OutOrStdout().Write([]byte(getVersion()))
	}

	return err
}

func getVersion() string {
	return Version + "\n"
}

func getFullVersion() string {
	vcsRevision := "unknown"
	vcsTime := "unknown"
	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				vcsRevision = setting.Value
			} else if setting.Key == "vcs.time" {
				vcsTime = setting.Value
			}
		}
	}

	return fmt.Sprintf(
		"%s (revision %s on %s)\n",
		Version,
		vcsRevision,
		vcsTime,
	)
}
