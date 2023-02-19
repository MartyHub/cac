package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	pv, found := os.LookupEnv(xdgConfigHome)
	_ = os.Setenv(xdgConfigHome, "/my_configs")
	t.Cleanup(func() {
		if found {
			_ = os.Setenv(xdgConfigHome, pv)
		} else {
			_ = os.Unsetenv(xdgConfigHome)
		}
	})

	configPath, err := GetConfigPath()

	assert.NoError(t, err)
	assert.Equal(t, "/my_configs/cac", configPath)
}
