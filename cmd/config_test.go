package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_runConfigList(t *testing.T) {
	const xdgConfigHome = "XDG_CONFIG_HOME"

	t.Setenv(xdgConfigHome, "../.config")

	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	require.NoError(t, runConfigList(cmd, false))
	assert.Equal(t, "json_config\n", buf.String())
}

func Test_readConfigAlias(t *testing.T) {
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
			if _, err := readConfigAlias(tt.args.configHome, tt.args.alias); (err != nil) != tt.wantErr {
				t.Errorf("loadConfigAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
