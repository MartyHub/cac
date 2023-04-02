package internal

import (
	"fmt"
	"strings"
	"time"
)

const (
	defaultTimeout  = 30 * time.Second
	defaultExpiry   = 11 * time.Hour
	defaultMaxConns = 4
	defaultMaxTries = 3
	defaultWait     = 100 * time.Millisecond
)

type Config struct {
	Aliases    []string      `json:"aliases"`
	AppID      string        `json:"app-id"`    //nolint:tagliatelle
	CertFile   string        `json:"cert-file"` //nolint:tagliatelle
	Expiry     time.Duration `json:"expiry"`
	Host       string        `json:"host"`
	KeyFile    string        `json:"key-file"`        //nolint:tagliatelle
	MaxConns   int           `json:"max-connections"` //nolint:tagliatelle
	MaxTries   int           `json:"max-tries"`       //nolint:tagliatelle
	Safe       string        `json:"safe"`
	SkipVerify bool          `json:"skip-verify"` //nolint:tagliatelle
	Timeout    time.Duration `json:"timeout"`
	Wait       time.Duration `json:"wait"`
}

func NewConfig() Config {
	return Config{
		Expiry:   defaultExpiry,
		MaxConns: defaultMaxConns,
		MaxTries: defaultMaxTries,
		Timeout:  defaultTimeout,
		Wait:     defaultWait,
	}
}

func (c Config) Overwrite(other Config) Config { //nolint:cyclop
	for _, alias := range other.Aliases {
		if !Contains(c.Aliases, alias) {
			c.Aliases = append(c.Aliases, alias)
		}
	}

	if other.AppID != "" {
		c.AppID = other.AppID
	}

	if other.CertFile != "" {
		c.CertFile = other.CertFile
	}

	if other.Expiry != defaultExpiry {
		c.Expiry = other.Expiry
	}

	if other.Host != "" {
		c.Host = other.Host
	}

	if other.KeyFile != "" {
		c.KeyFile = other.KeyFile
	}

	if other.MaxConns != defaultMaxConns {
		c.MaxConns = other.MaxConns
	}

	if other.MaxTries != defaultMaxTries {
		c.MaxTries = other.MaxTries
	}

	if other.Safe != "" {
		c.Safe = other.Safe
	}

	c.SkipVerify = other.SkipVerify

	if other.Timeout != defaultTimeout {
		c.Timeout = other.Timeout
	}

	if other.Wait != defaultWait {
		c.Wait = other.Wait
	}

	return c
}

func (c Config) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "aliases", strings.Join(c.Aliases, ", ")))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "app-id", c.AppID))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "cert-file", c.CertFile))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "expiry", c.Expiry))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "host", c.Host))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "key-file", c.KeyFile))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "max-conns", c.MaxConns))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "max-tries", c.MaxTries))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "safe", c.Safe))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "skip-verify", c.SkipVerify))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "timeout", c.Timeout))
	sb.WriteString(fmt.Sprintf("  %-11s = %v\n", "wait", c.Wait))

	return sb.String()
}
