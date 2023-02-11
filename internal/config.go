package internal

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

const xdgConfigHome = "XDG_CONFIG_HOME"

type Config struct {
	Name string `json:"name"`

	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`

	Host  string `json:"host"`
	AppId string `json:"appId"`
	Safe  string `json:"safe"`

	Json     *bool          `json:"json"`
	MaxConns *int           `json:"maxConns"`
	MaxTries *int           `json:"maxTries"`
	Timeout  *time.Duration `json:"timeout"`
}

func (c Config) Overwrite(params Parameters) Parameters {
	if params.CertFile == "" {
		params.CertFile = c.CertFile
	}

	if params.KeyFile == "" {
		params.KeyFile = c.KeyFile
	}

	if params.Host == "" {
		params.Host = c.Host
	}

	if params.AppId == "" {
		params.AppId = c.AppId
	}

	if params.Safe == "" {
		params.Safe = c.Safe
	}

	if !params.providedFlags["json"] && c.Json != nil {
		params.Json = *c.Json
	}

	if !params.providedFlags["maxConns"] && c.MaxConns != nil {
		params.MaxConns = *c.MaxConns
	}

	if !params.providedFlags["maxTries"] && c.MaxTries != nil {
		params.MaxTries = *c.MaxTries
	}

	if !params.providedFlags["timeout"] && c.Timeout != nil {
		params.Timeout = *c.Timeout
	}

	return params
}

func getConfigFile(params Parameters) string {
	configHome := os.Getenv(xdgConfigHome)

	if configHome == "" {
		userHome, err := os.UserHomeDir()

		if err != nil {
			params.Fatalf("Failed to get user home dir: %v", err)
		}

		configHome = path.Join(userHome, ".config")
	}

	return path.Join(configHome, "cac", "config.json")
}

func readConfigs(params Parameters) []Config {
	configFile := getConfigFile(params)
	bytes, err := os.ReadFile(configFile)

	var result []Config

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return result
		} else {
			params.Fatalf("Failed to read config file %s: %v", configFile, err)
		}
	}

	if err := json.Unmarshal(bytes, &result); err != nil {
		params.Fatalf("Failed to parse config file %s: %v", configFile, err)
	}

	return result
}

func GetConfig(params Parameters) Config {
	if params.Config != "" {
		for _, conf := range readConfigs(params) {
			if conf.Name == params.Config {
				return conf
			}
		}

		params.Fatalf("Failed to find config %s in %s", params.Config, getConfigFile(params))
	}

	return Config{}
}
