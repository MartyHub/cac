package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigHome(t *testing.T) {
	pv, found := os.LookupEnv(xdgConfigHome)
	_ = os.Setenv(xdgConfigHome, "../.config")
	t.Cleanup(func() {
		if found {
			_ = os.Setenv(xdgConfigHome, pv)
		} else {
			_ = os.Unsetenv(xdgConfigHome)
		}
	})

	home, err := GetConfigHome()

	assert.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}

func TestGetStateHome(t *testing.T) {
	pv, found := os.LookupEnv(xdgStateHome)
	_ = os.Setenv(xdgStateHome, "../.config")
	t.Cleanup(func() {
		if found {
			_ = os.Setenv(xdgStateHome, pv)
		} else {
			_ = os.Unsetenv(xdgStateHome)
		}
	})

	home, err := GetStateHome()

	assert.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}
