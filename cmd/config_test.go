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

	assert.NoError(t, runConfigList(cmd, false))
	assert.Equal(t, "json_config\nyaml_config\n", buf.String())
}

func Test_loadConfigAlias(t *testing.T) {
	type args struct {
		configHome string
		alias      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unknown alias",
			args: args{
				configHome: "../.config/cac",
				alias:      "unknown",
			},
			wantErr: true,
		},
		{
			name: "valid alias",
			args: args{
				configHome: "../.config/cac",
				alias:      "a1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadConfigAlias(tt.args.configHome, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("loadConfigAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
