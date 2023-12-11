package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigHome(t *testing.T) {
	t.Setenv(xdgConfigHome, "../.config")

	home, err := GetConfigHome()

	require.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}

func TestGetStateHome(t *testing.T) {
	t.Setenv(xdgStateHome, "../.config")

	home, err := GetStateHome()

	require.NoError(t, err)
	assert.Equal(t, "../.config/cac", home)
}
