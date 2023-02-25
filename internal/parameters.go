package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type Parameters struct {
	Config string

	CertFile string `mapstructure:"cert-file"`
	KeyFile  string `mapstructure:"key-file"`

	Host    string
	AppId   string `mapstructure:"app-id"`
	Safe    string
	Objects []string

	Json     bool
	MaxConns int `mapstructure:"max-connections"`
	MaxTries int `mapstructure:"max-tries"`
	Expiry   time.Duration
	Timeout  time.Duration
	Wait     time.Duration

	log *log.Logger
}

func NewParameters() Parameters {
	return Parameters{
		log: log.New(os.Stderr, "", 0),
	}
}

func (p Parameters) Errorf(format string, v ...any) {
	p.log.Printf(format, v...)
}

func (p Parameters) Fatalf(format string, v ...any) {
	p.log.Printf(format, v...)
	os.Exit(1)
}

func (p Parameters) Validate() error {
	errors := make([]string, 0)

	if p.CertFile == "" {
		errors = append(errors, "Certificate file is mandatory")
	}

	if p.KeyFile == "" {
		errors = append(errors, "Key file is mandatory")
	}

	if p.Host == "" {
		errors = append(errors, "Host is mandatory")
	}

	if p.AppId == "" {
		errors = append(errors, "Application Id is mandatory")
	}

	if p.Safe == "" {
		errors = append(errors, "Safe is mandatory")
	}

	if len(p.Objects) == 0 {
		if p.MaxConns <= 0 {
			errors = append(errors, "Either one object is required or max connections must be > 0")
		} else {
			stat, err := os.Stdin.Stat()

			if err != nil {
				p.Fatalf("Failed to stat stdin: %v", err)
			}

			if stat.Mode()&os.ModeCharDevice != 0 {
				errors = append(errors, "Either one object or a pipe is required")
			}
		}
	}

	if p.MaxConns < 0 {
		errors = append(errors, fmt.Sprintf("Max connections must be >= 0: %v", p.MaxConns))
	}

	if p.MaxTries <= 0 {
		errors = append(errors, fmt.Sprintf("Max tries must be > 0: %v", p.MaxTries))
	}

	if len(errors) > 0 {
		p.Errorf(strings.Join(errors, "\n"))
		return pflag.ErrHelp
	}

	return nil
}

func (p Parameters) fromStdin() bool {
	return len(p.Objects) == 0
}
