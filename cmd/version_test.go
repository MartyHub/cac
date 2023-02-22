package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getFullVersion(t *testing.T) {
	assert.Equal(t, "unknown (revision unknown on unknown)", getFullVersion())
}
