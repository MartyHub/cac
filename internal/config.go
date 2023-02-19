package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
)

const (
	xdgConfigHome              = "XDG_CONFIG_HOME"
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

func GetConfigPath() (string, error) {
	configHome, found := os.LookupEnv(xdgConfigHome)

	if !found {
		userHome, err := os.UserHomeDir()

		if err != nil {
			return "", err
		}

		configHome = path.Join(userHome, ".config")
	}

	return path.Join(configHome, "cac"), nil
}
