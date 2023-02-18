package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toPointer[T any](value T) *T {
	return &value
}

func TestConfig_Overwrite(t *testing.T) {
	config := Config{
		Aliases:  []string{"testConfig"},
		CertFile: "configCertFile",
		KeyFile:  "configKeyFile",
		MaxConns: toPointer(8),
		MaxTries: toPointer(4),
	}
	params := config.Overwrite(Parameters{
		CertFile: "certFile",
		MaxConns: 4,
		MaxTries: 8,
		providedFlags: map[string]bool{
			"maxTries": true,
		},
	})

	assert.Equal(t, "certFile", params.CertFile)
	assert.Equal(t, "configKeyFile", params.KeyFile)
	assert.Equal(t, 8, params.MaxConns)
	assert.Equal(t, 8, params.MaxTries)
}

func TestGetConfig(t *testing.T) {
	_ = os.Setenv(xdgConfigHome, "./.config")
	t.Cleanup(func() {
		_ = os.Unsetenv(xdgConfigHome)
	})

	config := GetConfig(Parameters{
		Config: "test_config",
	})

	assert.Equal(t, []string{"test_config"}, config.Aliases)
	assert.Equal(t, "localhost", config.Host)
}

func Test_getConfigFile(t *testing.T) {
	_ = os.Setenv(xdgConfigHome, "/my_configs")
	t.Cleanup(func() {
		_ = os.Unsetenv(xdgConfigHome)
	})

	assert.Equal(t, "/my_configs/cac/config.json", getConfigFile(newParameters()))
}
