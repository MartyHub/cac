package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestParameters(t *testing.T, ts *httptest.Server) Parameters {
	t.Helper()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)

	result := NewParameters()

	result.AppID = "appId"
	result.Host = u.Host
	result.MaxTries = 2
	result.Objects = []string{"o1", "o2"}
	result.Safe = "safe"

	return result
}

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	result := httptest.NewTLSServer(handler)

	t.Cleanup(func() {
		result.Close()
	})

	return result
}

var now = time.Unix(1677008748, 0).UTC() //nolint:gochecknoglobals

func newFixedClock() fixedClock {
	return fixedClock{
		t: now,
	}
}

func newTestClient(t *testing.T, handler http.HandlerFunc) Client {
	t.Helper()

	ts := newTestServer(t, handler)

	return Client{
		cache:  &Cache{},
		clock:  newFixedClock(),
		http:   ts.Client(),
		log:    log.New(io.Discard, "", 0),
		params: newTestParameters(t, ts),
	}
}

func captureOutput(client Client) *bytes.Buffer {
	result := &bytes.Buffer{}

	client.log.SetOutput(result)

	return result
}

func TestClient_Run(t *testing.T) {
	client := newTestClient(
		t,
		func(w http.ResponseWriter, r *http.Request) {
			object := r.URL.Query().Get("Object")
			_, _ = fmt.Fprintf(w, "{\"Content\": \"value for %s\"}\n", object)
		},
	)
	buf := captureOutput(client)

	assert.NoError(t, client.Run())
	assert.Equal(t, "o1='value for o1'\no2='value for o2'\n", buf.String())
}

func TestClient_Run_JSON(t *testing.T) {
	client := newTestClient(
		t,
		func(w http.ResponseWriter, r *http.Request) {
			object := r.URL.Query().Get("Object")
			_, _ = fmt.Fprintf(w, "{\"Content\": \"value for %s\"}\n", object)
		},
	)
	client.params.JSON = true
	buf := captureOutput(client)

	assert.NoError(t, client.Run())
	assert.Equal(
		t,
		`[
  {
    "object": "o1",
    "value": "value for o1",
    "try": 1,
    "statusCode": 200,
    "timestamp": "2023-02-21T19:45:48Z"
  },
  {
    "object": "o2",
    "value": "value for o2",
    "try": 1,
    "statusCode": 200,
    "timestamp": "2023-02-21T19:45:48Z"
  }
]
`,
		buf.String(),
	)
}

func TestClient_Run_BadRequest(t *testing.T) {
	client := newTestClient(
		t,
		func(w http.ResponseWriter, r *http.Request) {
			object := r.URL.Query().Get("Object")

			if object == "o1" {
				_, _ = fmt.Fprintf(w, "{\"Content\": \"value for %s\"}\n", object)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		},
	)
	buf := captureOutput(client)

	assert.Error(t, client.Run())
	assert.Equal(t, "o1='value for o1'\n", buf.String())
}

func TestClient_Run_InvalidJson(t *testing.T) {
	client := newTestClient(
		t,
		func(w http.ResponseWriter, r *http.Request) {
			object := r.URL.Query().Get("Object")

			if object == "o1" {
				_, _ = fmt.Fprintf(w, "{\"Content\": \"value for %s\"}\n", object)
			} else {
				_, _ = fmt.Fprintln(w, "Invalid JSON")
			}
		},
	)
	buf := captureOutput(client)

	assert.Error(t, client.Run())
	assert.Equal(t, "o1='value for o1'\n", buf.String())
}

func TestClient_Run_Retry(t *testing.T) {
	mu := sync.RWMutex{}
	objects := make(map[string]bool, 2)
	client := newTestClient(
		t,
		func(w http.ResponseWriter, r *http.Request) {
			object := r.URL.Query().Get("Object")

			mu.Lock()
			defer mu.Unlock()

			if objects[object] {
				_, _ = fmt.Fprintf(w, "{\"Content\": \"value for %s\"}\n", object)
			} else {
				objects[object] = true
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		},
	)
	buf := captureOutput(client)

	assert.NoError(t, client.Run())
	assert.Equal(t, "o1='value for o1'\no2='value for o2'\n", buf.String())
	assert.Contains(t, objects, "o1")
	assert.Contains(t, objects, "o2")
}

func TestClient_poolSize(t *testing.T) {
	tests := []struct {
		name   string
		client Client
		want   int
	}{
		{
			name: "maxConns",
			client: Client{
				params: Parameters{
					MaxConns: 2,
					Objects:  []string{"o1", "o2", "o3", "o4"},
				},
			},
			want: 2,
		},
		{
			name: "maxConns > length",
			client: Client{
				params: Parameters{
					MaxConns: 4,
					Objects:  []string{"o1", "o2"},
				},
			},
			want: 2,
		},
		{
			name: "length",
			client: Client{
				params: Parameters{
					Objects: []string{"o1", "o2", "o3", "o4"},
				},
			},
			want: 4,
		},
		{
			name: "stdin",
			client: Client{
				params: Parameters{
					MaxConns: 2,
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.client.poolSize())
		})
	}
}

func TestClient_query(t *testing.T) {
	client := &Client{
		params: Parameters{
			AppID: "appId",
			Safe:  "safe",
		},
	}

	assert.Equal(
		t,
		url.Values{
			"AppID":  []string{"appId"},
			"Safe":   []string{"safe"},
			"Object": []string{"o1"},
		},
		client.query("o1"),
	)
}

func TestClient_url(t *testing.T) {
	client := &Client{
		params: Parameters{
			Host: "host",
		},
	}

	assert.Equal(
		t,
		&url.URL{
			Scheme:   "https",
			Host:     "host",
			Path:     "/AIMWebService/api/Accounts",
			RawQuery: "AppID=appId&Safe=safe",
		},
		client.url(
			url.Values{
				"AppID": []string{"appId"},
				"Safe":  []string{"safe"},
			},
		),
	)
}

func TestClient_lineRegex(t *testing.T) {
	assert.Nil(t, lineRegex.FindStringSubmatch("KEY=VALUE"))
	assert.Equal(
		t,
		[]string{"KEY=${CYBERARK:OBJECT}", "KEY", "", "OBJECT", ""},
		lineRegex.FindStringSubmatch("KEY=${CYBERARK:OBJECT}"),
	)
	assert.Equal(
		t,
		[]string{"KEY=PREFIX_${CYBERARK:OBJECT}_SUFFIX", "KEY", "PREFIX_", "OBJECT", "_SUFFIX"},
		lineRegex.FindStringSubmatch("KEY=PREFIX_${CYBERARK:OBJECT}_SUFFIX"),
	)
}

func TestClient_readFromReader(t *testing.T) {
	client := Client{
		clock: newFixedClock(),
		log:   log.New(io.Discard, "", 0),
	}
	buf := captureOutput(client)
	in := make(chan *account, 1)

	assert.Equal(
		t,
		1,
		client.readFromReader(in, strings.NewReader("KEY1=${CYBERARK:o1}\nKEY2=VALUE2")),
	)

	assert.Equal(t, "KEY2=VALUE2\n", buf.String())

	result := <-in

	assert.Equal(t, "o1", result.Object)
	assert.Equal(t, "KEY1", result.key)
	assert.Equal(t, "", result.prefix)
	assert.Equal(t, "", result.suffix)
}
