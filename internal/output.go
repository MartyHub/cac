package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

func shellOutput(accounts []account) string {
	sb := strings.Builder{}

	for i, acct := range accounts {
		if acct.ok() {
			if i > 0 {
				sb.WriteRune('\n')
			}

			sb.WriteString(fmt.Sprintf("%s='%s'", acct.Object, acct.Value))
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
