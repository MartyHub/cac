package internal

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewParameters(t *testing.T) {
	p := NewParameters()

	assert.NotNil(t, p.log)
}

//nolint:funlen
func TestParameters_Valid(t *testing.T) {
	tests := []struct {
		name    string
		params  Parameters
		wantErr bool
	}{
		{
			name: "valid",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					AppID:    "appId",
					CertFile: "certFile",
					Host:     "host",
					KeyFile:  "keyFile",
					MaxTries: 1,
					Safe:     "safe",
				},
				Objects: []string{"object1"},
				JSON:    false,
			},
		},
		{
			name: "certFile",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					AppID:    "appId",
					Host:     "host",
					KeyFile:  "keyFile",
					MaxTries: 1,
					Safe:     "safe",
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "keyFile",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					AppID:    "appId",
					CertFile: "certFile",
					Host:     "host",
					MaxTries: 1,
					Safe:     "safe",
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "host",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					AppID:    "appId",
					CertFile: "certFile",
					KeyFile:  "keyFile",
					MaxTries: 1,
					Safe:     "safe",
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "appId",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					CertFile: "certFile",
					Host:     "host",
					KeyFile:  "keyFile",
					MaxTries: 1,
					Safe:     "safe",
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "safe",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					AppID:    "appId",
					CertFile: "certFile",
					KeyFile:  "keyFile",
					Host:     "host",
					MaxTries: 1,
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "objects",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					CertFile: "certFile",
					KeyFile:  "keyFile",
					Host:     "host",
					AppID:    "appId",
					Safe:     "safe",
					MaxConns: 0,
					MaxTries: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "maxConns",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					CertFile: "certFile",
					KeyFile:  "keyFile",
					Host:     "host",
					AppID:    "appId",
					Safe:     "safe",
					MaxConns: -1,
					MaxTries: 1,
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
		{
			name: "maxConns without object",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					CertFile: "certFile",
					KeyFile:  "keyFile",
					Host:     "host",
					AppID:    "appId",
					Safe:     "safe",
					MaxTries: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "maxTries",
			params: Parameters{
				log: log.New(os.Stderr, "", 0),
				Config: Config{
					CertFile: "certFile",
					KeyFile:  "keyFile",
					Host:     "host",
					AppID:    "appId",
					Safe:     "safe",
				},
				Objects: []string{"object1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if tt.wantErr {
				assert.ErrorIs(t, err, pflag.ErrHelp)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
