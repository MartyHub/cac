package internal

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewParameters(t *testing.T) {
	p := newParameters()

	assert.NotNil(t, p.log)
}

func TestParameters_Valid(t *testing.T) {
	tests := []struct {
		name        string
		params      Parameters
		wantValid   bool
		wantMessage string
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
			wantValid:   true,
			wantMessage: "",
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
			wantValid:   false,
			wantMessage: "Certificate file is mandatory\n",
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
			wantValid:   false,
			wantMessage: "Key file is mandatory\n",
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
			wantValid:   false,
			wantMessage: "Host is mandatory\n",
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
			wantValid:   false,
			wantMessage: "Application Id is mandatory\n",
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
			wantValid:   false,
			wantMessage: "Safe is mandatory\n",
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
			wantValid:   false,
			wantMessage: "Either one object is required or max conns must be > 0\n",
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
			wantValid:   false,
			wantMessage: "Max conns must be >= 0: -1\n",
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
			wantValid:   false,
			wantMessage: "Either one object is required or max conns must be > 0\n",
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
			wantValid:   false,
			wantMessage: "Max tries must be > 0: 0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, message := tt.params.Valid()
			assert.Equal(t, tt.wantValid, valid)
			assert.Equal(t, tt.wantMessage, message)
		})
	}
}

func Test_parse(t *testing.T) {
	args := []string{
		"cac",
		"-certFile", "certFile",
		"-keyFile", "keyFile",
		"-host", "host",
		"-appId", "appId",
		"-safe", "safe",
		"-object", "o1",
		"-object", "o2",
	}

	params := Parse(args)

	assert.Equal(t, "certFile", params.CertFile)
	assert.Equal(t, "keyFile", params.KeyFile)
	assert.Equal(t, "host", params.Host)
	assert.Equal(t, "appId", params.AppId)
	assert.Equal(t, "safe", params.Safe)
	assert.Equal(t, []string{"o1", "o2"}, params.Objects)

	// Default values
	assert.False(t, params.Json)
	assert.Equal(t, 4, params.MaxConns)
	assert.Equal(t, 3, params.MaxTries)
	assert.Equal(t, 30*time.Second, params.Timeout)
	assert.Equal(t, 100*time.Millisecond, params.Wait)

	assert.NotContains(t, params.providedFlags, "json")
	assert.NotContains(t, params.providedFlags, "maxConns")
	assert.NotContains(t, params.providedFlags, "maxTries")
	assert.NotContains(t, params.providedFlags, "maxTimeout")
}

func TestParameters_getVersion(t *testing.T) {
	assert.Equal(t, "unknown (revision unknown on unknown)", newParameters().getVersion())
}
