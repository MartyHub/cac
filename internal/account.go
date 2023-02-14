package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type successBody struct {
	Content string `json:"Content"`
}

type errorBody struct {
	ErrorCode string `json:"ErrorCode"`
	ErrorMsg  string `json:"ErrorMsg"`
}

type account struct {
	Object         string `json:"object"`
	Value          string `json:"value"`
	Try            int    `json:"try"`
	Error          error  `json:"error,omitempty"`
	StatusCode     int    `json:"statusCode"`
	prefix, suffix string
}

func newAccount(object string, prefix, suffix string) *account {
	return &account{
		Object: object,
		prefix: prefix,
		suffix: suffix,
	}
}

func (acct *account) newTry() {
	acct.Try++
	acct.Error = nil
	acct.StatusCode = 0
}

func (acct *account) retry(maxTries int) bool {
	if acct.Try >= maxTries {
		return false
	}

	return acct.StatusCode == 0 ||
		acct.StatusCode == http.StatusInternalServerError ||
		acct.StatusCode == http.StatusBadGateway ||
		acct.StatusCode == http.StatusServiceUnavailable ||
		acct.StatusCode == http.StatusGatewayTimeout
}

func (acct *account) ok() bool {
	return acct.Error == nil && acct.StatusCode == http.StatusOK
}

func (acct *account) parseError(data []byte) {
	var result *errorBody

	if err := parseBody(data, &result); err != nil {
		acct.Error = fmt.Errorf("failed to parse JSON '%s'", string(data))
	} else {
		acct.Error = fmt.Errorf("%s: %s", result.ErrorCode, result.ErrorMsg)
	}
}

func (acct *account) parseSuccess(data []byte) {
	var result *successBody

	if err := parseBody(data, &result); err != nil {
		acct.Error = fmt.Errorf("failed to parse JSON '%s'", string(data))
	} else {
		acct.Value = result.Content
	}
}

func (acct *account) shell(fromStdin bool) string {
	if fromStdin {
		return strings.Join([]string{acct.prefix, acct.Value, acct.suffix}, "")
	} else {
		return fmt.Sprintf("%s='%s'", acct.Object, acct.Value)
	}
}

func (acct *account) String() string {
	return fmt.Sprintf("%s # %d: status=%d, error=%v", acct.Object, acct.Try, acct.StatusCode, acct.Error)
}

func parseBody[T any](data []byte, result *T) error {
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	return nil
}
