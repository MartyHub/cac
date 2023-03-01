package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func fileOutput(accounts []account, output string) error {
	err := os.MkdirAll(output, 0o700)
	if err != nil {
		return err
	}

	for _, acct := range accounts {
		err := os.WriteFile(filepath.Join(output, acct.key), []byte(acct.Value), 0o600)
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
