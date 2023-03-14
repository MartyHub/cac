package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigHome(t *testing.T) {
	t.Setenv(xdgConfigHome, "../.config")

	home, err := GetConfigHome()

	assert.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}

func TestGetStateHome(t *testing.T) {
	t.Setenv(xdgStateHome, "../.config")

	home, err := GetStateHome()

	assert.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}
