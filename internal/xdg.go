package internal

import (
	"os"
	"path"
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

		home = path.Join(userHome, ".config")
	}

	result := path.Join(home, "cac")

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

		home = path.Join(userHome, ".local", "state")
	}

	result := path.Join(home, "cac")

	if err := os.MkdirAll(result, rwx); err != nil {
		return "", err
	}

	return result, nil
}
