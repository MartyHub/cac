package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var accounts = []account{
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
		Error:     fmt.Errorf("test error"),
		Timestamp: time.Unix(1677008748, 0).UTC(),
	},
}

func Test_jsonOutput(t *testing.T) {
	output, err := jsonOutput(accounts)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"[\n  {\n    \"object\": \"object1\",\n    \"value\": \"value1\",\n    \"try\": 1,\n    \"statusCode\": 200,\n    \"timestamp\": \"2023-02-21T19:45:48Z\"\n  },\n  {\n    \"object\": \"object\",\n    \"value\": \"\",\n    \"try\": 3,\n    \"error\": {},\n    \"statusCode\": 0,\n    \"timestamp\": \"2023-02-21T19:45:48Z\"\n  }\n]",
		output,
	)
}

func Test_shellOutput(t *testing.T) {
	assert.Equal(t, "object1='value1'", shellOutput(accounts, false))
}
