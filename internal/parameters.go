package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Parameters struct {
	Config

	CfgName string
	JSON    bool
	Objects []string
	Output  string

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

//nolint:cyclop
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

	if p.AppID == "" {
		errors = append(errors, "Application Id is mandatory")
	}

	if p.Safe == "" {
		errors = append(errors, "Safe is mandatory")
	}

	if len(p.Objects) == 0 {
		errors = p.validateObject(errors)
	}

	if p.MaxConns < 0 {
		errors = append(errors, fmt.Sprintf("Max connections must be >= 0: %v", p.MaxConns))
	}

	if p.MaxTries <= 0 {
		errors = append(errors, fmt.Sprintf("Max tries must be > 0: %v", p.MaxTries))
	}

	if p.Output != "" && !p.fromStdin() {
		errors = append(errors, "no args should be given if output is set")
	}

	if len(errors) > 0 {
		p.Errorf(strings.Join(errors, "\n"))

		return pflag.ErrHelp
	}

	return nil
}

func (p Parameters) validateObject(errors []string) []string {
	if p.MaxConns <= 0 {
		return append(errors, "Either one object is required or max connections must be > 0")
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		p.Fatalf("Failed to stat stdin: %v", err)
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		errors = append(errors, "Either one object or a pipe is required")
	}

	return errors
}

func (p Parameters) fromStdin() bool {
	return len(p.Objects) == 0
}
