package internal

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type Client struct {
	http   *http.Client
	log    *log.Logger // help testing
	params Parameters
}

func NewClient(params Parameters) Client {
	cert, err := tls.LoadX509KeyPair(params.CertFile, params.KeyFile)

	if err != nil {
		params.Fatalf("Failed to read certificate %s and key %s: %v", params.CertFile, params.KeyFile, err)
	}

	return Client{
		http: &http.Client{
			Timeout: params.Timeout,
			Transport: &http.Transport{
				MaxConnsPerHost: params.MaxConns,
				Proxy:           nil,
				TLSClientConfig: &tls.Config{
					Certificates:  []tls.Certificate{cert},
					MinVersion:    tls.VersionTLS12,
					Renegotiation: tls.RenegotiateOnceAsClient,
				},
			},
		},
		log:    log.New(os.Stdout, "", 0),
		params: params,
	}
}

func (c Client) Run() bool {
	l := c.length()
	in := make(chan *account, l)
	out := make(chan *account, l)

	for i := 0; i < c.poolSize(); i++ {
		go c.worker(in, out)
	}

	for _, object := range c.params.Objects {
		in <- newAccount(object)
	}

	close(in)

	accounts := c.collect(out)

	if c.params.Json {
		output, err := jsonOutput(accounts)

		if err != nil {
			c.params.Fatalf("Failed to marshall result as JSON: %v", err)
		}

		c.log.Print(output)
	} else {
		c.log.Print(shellOutput(accounts))
	}

	return c.ok(accounts)
}

func (c Client) length() int {
	return len(c.params.Objects)
}

func (c Client) poolSize() int {
	l := c.length()

	if c.params.MaxConns == 0 || l < c.params.MaxConns {
		return l
	}

	return c.params.MaxConns
}

func (c Client) collect(accounts <-chan *account) []account {
	l := c.length()
	results := make([]account, 0, l)

	for i := 0; i < l; i++ {
		acct := <-accounts

		results = append(results, *acct)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Object < results[j].Object
	})

	return results
}

func (c Client) ok(accounts []account) bool {
	for _, acct := range accounts {
		if !acct.ok() {
			return false
		}
	}

	return true
}

func (c Client) worker(in <-chan *account, out chan<- *account) {
	for acct := range in {
		for acct.run(c.params.MaxTries) {
			time.Sleep(time.Duration(acct.Try*acct.Try*100) * time.Millisecond)
			acct.newTry()
			c.get(acct)
		}

		out <- acct
	}
}

func (c Client) url(values url.Values) *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     c.params.Host,
		Path:     "/AIMWebService/api/Accounts",
		RawQuery: values.Encode(),
	}
}

func (c Client) query(object string) url.Values {
	result := url.Values{}

	result.Set("AppID", c.params.AppId)
	result.Set("Safe", c.params.Safe)
	result.Set("Object", object)

	return result
}

func (c Client) get(acct *account) {
	response, err := c.http.Get(c.url(c.query(acct.Object)).String())

	if err != nil {
		acct.Error = err
		return
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			c.params.Errorf("Failed to close body: %v", err)
		}
	}()

	acct.StatusCode = response.StatusCode
	data, err := io.ReadAll(response.Body)

	if err != nil {
		acct.Error = err
		return
	}

	if response.StatusCode == http.StatusOK {
		acct.parseSuccess(data)
	} else {
		acct.parseError(data)
	}
}
