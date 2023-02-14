package internal

import (
	"encoding/json"
	"strings"
)

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

func jsonOutput(accounts []account) (string, error) {
	bytes, err := json.MarshalIndent(accounts, "", "  ")

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
