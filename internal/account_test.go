package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_account_newTry(t *testing.T) {
	acct := &account{
		Object:     "object",
		Try:        1,
		Error:      fmt.Errorf("error"),
		StatusCode: 500,
	}

	acct.newTry()

	assert.Equal(t, "object", acct.Object)
	assert.Equal(t, 2, acct.Try)
	assert.Nil(t, acct.Error)
	assert.Zero(t, acct.StatusCode)
}

func Test_account_ok(t *testing.T) {
	tests := []struct {
		name string
		acct *account
		want bool
	}{
		{
			name: "err",
			acct: &account{
				Error: fmt.Errorf("test"),
			},
			want: false,
		},
		{
			name: "statusCode",
			acct: &account{
				StatusCode: 500,
			},
			want: false,
		},
		{
			name: "ok",
			acct: &account{
				StatusCode: 200,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.acct.ok())
		})
	}
}

func Test_account_parseError(t *testing.T) {
	type args struct {
		data []byte
	}

	tests := []struct {
		name string
		acct *account
		args args
		want error
	}{
		{
			name: "invalidJSON",
			acct: &account{},
			args: args{
				data: []byte("Invalid JSON"),
			},
			want: fmt.Errorf("failed to parse JSON 'Invalid JSON'"),
		},
		{
			name: "cyberArkError",
			acct: &account{},
			args: args{
				data: []byte("{\"ErrorCode\": \"code\", \"ErrorMsg\": \"message\"}"),
			},
			want: fmt.Errorf("code: message"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.acct.parseError(tt.args.data)
			assert.Equal(t, tt.acct.Error, tt.want)
		})
	}
}

func Test_account_parseSuccess(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name      string
		acct      *account
		args      args
		wantErr   error
		wantValue string
	}{
		{
			name: "invalidJSON",
			acct: &account{},
			args: args{
				data: []byte("Invalid JSON"),
			},
			wantErr: fmt.Errorf("failed to parse JSON 'Invalid JSON'"),
		},
		{
			name: "ok",
			acct: &account{},
			args: args{
				data: []byte(`{"Content": "value"}`),
			},
			wantValue: "value",
		},
		{
			name: "single quote",
			acct: &account{},
			args: args{
				data: []byte(`{"Content": "'value'"}`),
			},
			wantValue: "value",
		},
		{
			name: "double quote",
			acct: &account{},
			args: args{
				data: []byte(`{"Content": "\"value\""}`),
			},
			wantValue: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.acct.parseSuccess(tt.args.data)
			assert.Equal(t, tt.wantErr, tt.acct.Error)
			assert.Equal(t, tt.wantValue, tt.acct.Value)
		})
	}
}

func Test_account_retry(t *testing.T) {
	tests := []struct {
		name string
		acct *account
		want bool
	}{
		{
			name: "200",
			acct: &account{
				StatusCode: 200,
				Try:        1,
			},
			want: false,
		},
		{
			name: "500",
			acct: &account{
				StatusCode: 500,
				Try:        1,
			},
			want: true,
		},
		{
			name: "502",
			acct: &account{
				StatusCode: 502,
				Try:        2,
			},
			want: true,
		},
		{
			name: "503",
			acct: &account{
				StatusCode: 503,
				Try:        3,
			},
			want: true,
		},
		{
			name: "504",
			acct: &account{
				StatusCode: 504,
				Try:        4,
			},
			want: true,
		},
		{
			name: "504",
			acct: &account{
				StatusCode: 504,
				Try:        5,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.acct.retry(5))
		})
	}
}

func Test_newAccount(t *testing.T) {
	assert.Equal(
		t,
		&account{Object: "object", Timestamp: now},
		newAccount("object", now, "", "", ""),
	)
}

func Test_parseBody(t *testing.T) {
	type args[T any] struct {
		data   []byte
		result *T
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		wantErr bool
	}
	tests := []testCase[successBody]{
		{
			name: "error",
			args: args[successBody]{
				data:   []byte("Invalid JSON"),
				result: &successBody{},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args[successBody]{
				data:   []byte("{\"Content\": \"value\"}"),
				result: &successBody{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, parseBody(tt.args.data, tt.args.result) != nil)
		})
	}
}

func Test_account_shell(t *testing.T) {
	type args struct {
		fromStdin bool
	}
	tests := []struct {
		name string
		acct *account
		args args
		want string
	}{
		{
			name: "from params",
			acct: &account{
				Object: "object",
				Value:  "value",
			},
			args: args{fromStdin: false},
			want: "object='value'",
		},
		{
			name: "from stdin",
			acct: &account{
				Object: "object",
				Value:  "value",
				key:    "key",
				prefix: "prefix_",
				suffix: "_suffix",
			},
			args: args{fromStdin: true},
			want: "key=prefix_value_suffix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.acct.shell(tt.args.fromStdin), "shell(%v)", tt.args.fromStdin)
		})
	}
}
