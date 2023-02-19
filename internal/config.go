package internal

import (
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

const (
	configFilePerm fs.FileMode = 0o600
)

func CheckConfigFilePermissions(configFile string) error {
	if runtime.GOOS != "windows" {
		stat, err := os.Stat(configFile)

		if err != nil {
			return err
		}

		permissions := stat.Mode().Perm()

		if permissions != configFilePerm {
			if err = os.Chmod(configFile, configFilePerm); err != nil {
				return fmt.Errorf(
					"incorrect permissions %v for config file %s (must be %v): %w",
					permissions,
					configFile,
					configFilePerm,
					err,
				)
			}
		}
	}

	return nil
}
