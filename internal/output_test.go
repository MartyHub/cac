package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var accounts = []account{
	{
		Object:     "object1",
		Value:      "value1",
		Try:        1,
		Error:      nil,
		StatusCode: 200,
	},
	{
		Object: "object",
		Try:    3,
		Error:  fmt.Errorf("test error"),
	},
}

func Test_jsonOutput(t *testing.T) {
	output, err := jsonOutput(accounts)

	assert.NoError(t, err)
	assert.Equal(
		t,
		"[\n  {\n    \"object\": \"object1\",\n    \"value\": \"value1\",\n    \"try\": 1,\n    \"statusCode\": 200\n  },\n  {\n    \"object\": \"object\",\n    \"value\": \"\",\n    \"try\": 3,\n    \"error\": {},\n    \"statusCode\": 0\n  }\n]",
		output)
}

func Test_shellOutput(t *testing.T) {
	assert.Equal(t, "object1='value1'", shellOutput(accounts))
}
