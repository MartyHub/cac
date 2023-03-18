package internal

import (
	"os"
	"path/filepath"
)

const (
	xdgConfigHome = "XDG_CONFIG_HOME"
	xdgStateHome  = "XDG_STATE_HOME"
)

func GetConfigHome() (string, error) {
	home, found := os.LookupEnv(xdgConfigHome)

	if !found {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		home = filepath.Join(userHome, ".config")
	}

	result := filepath.Join(home, "cac")

	if err := os.MkdirAll(result, rwx); err != nil {
		return "", err
	}

	return result, nil
}

func GetStateHome() (string, error) {
	home, found := os.LookupEnv(xdgStateHome)

	if !found {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		home = filepath.Join(userHome, ".local", "state")
	}

	result := filepath.Join(home, "cac")

	if err := os.MkdirAll(result, rwx); err != nil {
		return "", err
	}

	return result, nil
}
