package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_runConfigList(t *testing.T) {
	const xdgConfigHome = "XDG_CONFIG_HOME"

	pv, found := os.LookupEnv(xdgConfigHome)
	_ = os.Setenv(xdgConfigHome, "../.config")
	t.Cleanup(func() {
		if found {
			_ = os.Setenv(xdgConfigHome, pv)
		} else {
			_ = os.Unsetenv(xdgConfigHome)
		}
	})

	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	assert.NoError(t, runConfigList(cmd))
	assert.Equal(t, "json_config\nyaml_config\n", buf.String())
}
