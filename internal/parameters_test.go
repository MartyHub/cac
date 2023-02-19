package internal

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParameters(t *testing.T) {
	p := NewParameters()

	assert.NotNil(t, p.log)
}

func TestParameters_Valid(t *testing.T) {
	tests := []struct {
		name   string
		params Parameters
		want   error
	}{
		{
			name: "valid",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				Objects:  []string{"object1"},
				Json:     false,
				MaxTries: 1,
			},
		},
		{
			name: "certFile",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				Objects:  []string{"object1"},
				MaxTries: 1,
			},
			want: fmt.Errorf("Certificate file is mandatory"),
		},
		{
			name: "keyFile",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				Objects:  []string{"object1"},
				MaxTries: 1,
			},
			want: fmt.Errorf("Key file is mandatory"),
		},
		{
			name: "host",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				AppId:    "appId",
				Safe:     "safe",
				Objects:  []string{"object1"},
				MaxTries: 1,
			},
			want: fmt.Errorf("Host is mandatory"),
		},
		{
			name: "appId",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				Safe:     "safe",
				Objects:  []string{"object1"},
				MaxTries: 1,
			},
			want: fmt.Errorf("Application Id is mandatory"),
		},
		{
			name: "safe",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Objects:  []string{"object1"},
				MaxTries: 1,
			},
			want: fmt.Errorf("Safe is mandatory"),
		},
		{
			name: "objects",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				MaxConns: 0,
				MaxTries: 1,
			},
			want: fmt.Errorf("Either one object is required or max connections must be > 0"),
		},
		{
			name: "maxConns",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				MaxConns: -1,
				MaxTries: 1,
				Objects:  []string{"object1"},
			},
			want: fmt.Errorf("Max connections must be >= 0: -1"),
		},
		{
			name: "maxConns without object",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				MaxTries: 1,
			},
			want: fmt.Errorf("Either one object is required or max connections must be > 0"),
		},
		{
			name: "maxTries",
			params: Parameters{
				log:      log.New(os.Stderr, "", 0),
				CertFile: "certFile",
				KeyFile:  "keyFile",
				Host:     "host",
				AppId:    "appId",
				Safe:     "safe",
				Objects:  []string{"object1"},
			},
			want: fmt.Errorf("Max tries must be > 0: 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			assert.Equal(t, tt.want, err)
		})
	}
}
