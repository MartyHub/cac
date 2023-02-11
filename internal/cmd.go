package internal

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

type Parameters struct {
	Config string

	CertFile string
	KeyFile  string

	Host    string
	AppId   string
	Safe    string
	Objects []string

	Json     bool
	MaxConns int
	MaxTries int
	Timeout  time.Duration

	log           *log.Logger
	providedFlags map[string]bool
}

func newParameters() Parameters {
	return Parameters{
		log:           log.New(os.Stderr, "", 0),
		providedFlags: make(map[string]bool),
	}
}

func (p Parameters) Errorf(format string, v ...any) {
	p.log.Printf(format, v...)
}

func (p Parameters) Fatalf(format string, v ...any) {
	p.log.Printf(format, v...)
	os.Exit(1)
}

func (p Parameters) Valid() bool {
	result := true

	if p.CertFile == "" {
		p.Errorf("Certificate file is mandatory")
		result = false
	}

	if p.KeyFile == "" {
		p.Errorf("Key file is mandatory")
		result = false
	}

	if p.Host == "" {
		p.Errorf("Host is mandatory")
		result = false
	}

	if p.AppId == "" {
		p.Errorf("Application Id is mandatory")
		result = false
	}

	if p.Safe == "" {
		p.Errorf("Safe is mandatory")
		result = false
	}

	if len(p.Objects) == 0 {
		p.Errorf("At least one object is mandatory")
		result = false
	}

	if p.MaxTries <= 0 {
		p.Errorf("Max tries must be > 0: %v", p.MaxTries)
		result = false
	}

	return result
}

func (p Parameters) provided(f *flag.Flag) {
	p.providedFlags[f.Name] = true
}

type stringsValue []string

func (f *stringsValue) String() string {
	return strings.Join(*f, ", ")
}

func (f *stringsValue) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func Parse(args []string) Parameters {
	var objects stringsValue

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	params := newParameters()

	flags.StringVar(&params.Config, "config", "", "Config name")

	flags.StringVar(&params.CertFile, "certFile", "", "Certificate file")
	flags.StringVar(&params.KeyFile, "keyFile", "", "Key file")

	flags.StringVar(&params.Host, "host", "", "CyberArk CCP REST Web Service Host")
	flags.StringVar(&params.AppId, "appId", "", "CyberArk Application Id")
	flags.StringVar(&params.Safe, "safe", "", "CyberArk Safe")
	flags.Var(&objects, "object", "CyberArk Object (at least one required)")

	flags.BoolVar(&params.Json, "json", false, "JSON output")
	flags.IntVar(&params.MaxConns, "maxConns", 4, "Max connections")
	flags.IntVar(&params.MaxTries, "maxTries", 3, "Max tries")
	flags.DurationVar(&params.Timeout, "timeout", 30*time.Second, "Timeout")

	// Ignore error as flags is set for ExitOnError
	_ = flags.Parse(args[1:])

	params.Objects = objects

	flags.Visit(params.provided)

	params = GetConfig(params).Overwrite(params)

	if !params.Valid() {
		flags.Usage()
		params.Fatalf("")
	}

	return params
}
