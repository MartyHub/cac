package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var accounts = []Account{ //nolint:gochecknoglobals
	{
		Object:     "object1",
		Value:      "value1",
		Try:        1,
		Error:      nil,
		StatusCode: 200,
		Timestamp:  time.Unix(1677008748, 0).UTC(),
	},
	{
		Object:    "object",
		Try:       3,
		Error:     NewError(nil, "test error"),
		Timestamp: time.Unix(1677008748, 0).UTC(),
	},
}

func Test_jsonOutput(t *testing.T) {
	output, err := jsonOutput(accounts)

	assert.NoError(t, err)
	assert.Equal(
		t,
		`[
  {
    "object": "object1",
    "value": "value1",
    "try": 1,
    "statusCode": 200,
    "timestamp": "2023-02-21T19:45:48Z"
  },
  {
    "object": "object",
    "value": "",
    "try": 3,
    "error": {
      "Cause": null,
      "Message": "test error"
    },
    "statusCode": 0,
    "timestamp": "2023-02-21T19:45:48Z"
  }
]`,
		output,
	)
}

func Test_shellOutput(t *testing.T) {
	assert.Equal(t, "object1='value1'", shellOutput(accounts, false))
}
