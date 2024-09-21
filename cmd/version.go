package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version is injected during the build.
var Version = "unknown" //nolint:gochecknoglobals

func newVersionCommand() *cobra.Command {
	full := false
	result := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Display version information",
		Run: func(cmd *cobra.Command, _ []string) {
			runVersion(cmd, full)
		},
	}

	result.Flags().BoolVarP(&full, "full", "f", false, "also display Git revision and modification time")

	return result
}

func runVersion(cmd *cobra.Command, full bool) {
	if full {
		cmd.Println(getFullVersion())
	} else {
		cmd.Println(Version)
	}
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
		"%s (revision %s on %s)",
		Version,
		vcsRevision,
		vcsTime,
	)
}
