package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var Version = "unknown"

type Parameters struct {
	Config string

	CertFile string
	KeyFile  string

	Host    string
	AppId   string
	Safe    string
	Objects []string

	Version bool

	Json     bool
	MaxConns int
	MaxTries int
	Timeout  time.Duration
	Wait     time.Duration

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

func (p Parameters) Valid() (bool, string) {
	sb := strings.Builder{}
	valid := true

	if p.CertFile == "" {
		sb.WriteString("Certificate file is mandatory\n")
		valid = false
	}

	if p.KeyFile == "" {
		sb.WriteString("Key file is mandatory\n")
		valid = false
	}

	if p.Host == "" {
		sb.WriteString("Host is mandatory\n")
		valid = false
	}

	if p.AppId == "" {
		sb.WriteString("Application Id is mandatory\n")
		valid = false
	}

	if p.Safe == "" {
		sb.WriteString("Safe is mandatory\n")
		valid = false
	}

	if len(p.Objects) == 0 {
		if p.MaxConns <= 0 {
			sb.WriteString("Either one object is required or max conns must be > 0\n")
			valid = false
		} else {
			stat, err := os.Stdin.Stat()

			if err != nil {
				p.Fatalf("Failed to stat stdin: %v", err)
			}

			if stat.Mode()&os.ModeCharDevice != 0 {
				sb.WriteString("Either one object or a pipe is required")
				valid = false
			}
		}
	}

	if p.MaxConns < 0 {
		sb.WriteString(fmt.Sprintf("Max conns must be >= 0: %v\n", p.MaxConns))
		valid = false
	}

	if p.MaxTries <= 0 {
		sb.WriteString(fmt.Sprintf("Max tries must be > 0: %v\n", p.MaxTries))
		valid = false
	}

	return valid, sb.String()
}

func (p Parameters) getVersion() string {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		p.Fatalf("Failed to read build info")
	}

	vcsRevision := "unknown"
	vcsTime := "unknown"

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			vcsRevision = setting.Value
		} else if setting.Key == "vcs.time" {
			vcsTime = setting.Value
		}
	}

	return fmt.Sprintf(
		"%s (revision %s on %s)",
		Version,
		vcsRevision,
		vcsTime,
	)
}

func (p Parameters) provided(f *flag.Flag) {
	p.providedFlags[f.Name] = true
}

func (p Parameters) fromStdin() bool {
	return len(p.Objects) == 0
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

	flags.StringVar(&params.Config, "c", "", "Config name")

	flags.StringVar(&params.CertFile, "certFile", "", "Certificate file")
	flags.StringVar(&params.KeyFile, "keyFile", "", "Key file")

	flags.StringVar(&params.Host, "host", "", "CyberArk CCP REST Web Service Host")
	flags.StringVar(&params.AppId, "appId", "", "CyberArk Application Id")
	flags.StringVar(&params.Safe, "safe", "", "CyberArk Safe")
	flags.Var(&objects, "o", "CyberArk Object")

	flags.BoolVar(&params.Json, "json", false, "JSON output")
	flags.IntVar(&params.MaxConns, "maxConns", 4, "Max connections")
	flags.IntVar(&params.MaxTries, "maxTries", 3, "Max tries")
	flags.DurationVar(&params.Timeout, "timeout", 30*time.Second, "Timeout")
	flags.DurationVar(&params.Wait, "wait", 100*time.Millisecond, "Wait before retry")

	flags.BoolVar(&params.Version, "version", false, "Display version information")

	// Ignore error as flags is set for ExitOnError
	_ = flags.Parse(args[1:])

	params.Objects = objects

	flags.Visit(params.provided)

	params = GetConfig(params).Overwrite(params)

	if params.Version {
		fmt.Println(params.getVersion())
		os.Exit(0)
	}

	valid, message := params.Valid()

	if !valid {
		flags.Usage()
		params.Fatalf(message)
	}

	return params
}
