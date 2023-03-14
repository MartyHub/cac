package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	rwx = 0o700
	rw  = 0o600
)

func fileOutput(accounts []account, output string) error {
	err := os.MkdirAll(output, rwx)
	if err != nil {
		return err
	}

	for _, acct := range accounts {
		err := os.WriteFile(filepath.Join(output, acct.key), []byte(acct.Value), rw)
		if err != nil {
			return err
		}
	}

	return nil
}

func jsonOutput(accounts []account) (string, error) {
	bytes, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func shellOutput(accounts []account, fromStdin bool) string {
	sb := strings.Builder{}

	for i, acct := range accounts {
		if acct.ok() {
			if i > 0 {
				sb.WriteRune('\n')
			}

			sb.WriteString(acct.shell(fromStdin))
		}
	}

	return sb.String()
}
