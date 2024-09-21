package internal

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"time"
)

var lineRegex = regexp.MustCompile(`^([^#][^=]*)=(.*)\$\{CYBERARK:(.+)}(.*)$`)

const (
	lineRegexGroups = 5
)

type Client struct {
	clock  clock
	http   *http.Client
	log    *log.Logger // help testing
	params Parameters
}

func NewClient(params Parameters) (Client, error) {
	cert, err := tls.LoadX509KeyPair(params.CertFile, params.KeyFile)
	if err != nil {
		return Client{}, err
	}

	return Client{
		clock: utcClock{},
		http: &http.Client{
			Timeout: params.Timeout,
			Transport: &http.Transport{
				MaxConnsPerHost: params.MaxConns,
				Proxy:           nil,
				TLSClientConfig: &tls.Config{
					Certificates:       []tls.Certificate{cert},
					MinVersion:         tls.VersionTLS12,
					Renegotiation:      tls.RenegotiateOnceAsClient,
					InsecureSkipVerify: params.SkipVerify, //nolint:gosec
				},
			},
		},
		log:    log.New(os.Stdout, "", 0),
		params: params,
	}, nil
}

func (c Client) Run() error {
	cache, err := NewDBCache()
	if err != nil {
		return err
	}

	defer cache.Close()

	if err = cache.clean(c.clock, c.params.Expiry); err != nil {
		return err
	}

	size := c.poolSize()
	in := make(chan *Account, size)
	out := make(chan *Account, size)

	for range size {
		go c.worker(cache, in, out)
	}

	count := make(chan int)

	go c.read(in, count)

	accounts := c.collect(out, count)

	close(in)

	if err = c.output(accounts); err != nil {
		return err
	}

	if err = cache.merge(c.params.CfgName, accounts); err != nil {
		return err
	}

	return c.ok(accounts)
}

func (c Client) output(accounts []Account) error {
	switch {
	case c.params.JSON:
		output, err := jsonOutput(accounts)
		if err != nil {
			return err
		}

		c.log.Print(output)
	case c.params.Output != "":
		return fileOutput(accounts, c.params.Output)
	default:
		c.log.Print(shellOutput(accounts, c.params.fromStdin()))
	}

	return nil
}

func (c Client) read(in chan<- *Account, count chan<- int) {
	if c.params.fromStdin() {
		count <- c.readFromReader(in, os.Stdin)
	} else {
		count <- c.readFromParams(in)
	}
}

func (c Client) readFromParams(in chan<- *Account) int {
	now := c.clock.now()
	result := 0

	for _, object := range c.params.Objects {
		in <- newAccount(object, now, "", "", "")

		result++
	}

	return result
}

func (c Client) readFromReader(in chan<- *Account, reader io.Reader) int {
	now := c.clock.now()
	scanner := bufio.NewScanner(reader)
	result := 0

	for scanner.Scan() {
		line := scanner.Text()
		groups := lineRegex.FindStringSubmatch(line)

		if len(groups) == lineRegexGroups {
			in <- newAccount(groups[3], now, groups[1], groups[2], groups[4])

			result++
		} else {
			c.log.Print(line)
		}
	}

	return result
}

func (c Client) poolSize() int {
	if !c.params.fromStdin() {
		l := len(c.params.Objects)

		if c.params.MaxConns == 0 || l < c.params.MaxConns {
			return l
		}
	}

	return c.params.MaxConns
}

func (c Client) collect(accounts <-chan *Account, count <-chan int) []Account {
	l := 0
	lenKnown := false
	results := make([]Account, 0)

	for !lenKnown || len(results) != l {
		select {
		case acct := <-accounts:
			results = append(results, *acct)
		case l = <-count:
			lenKnown = true
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Object < results[j].Object
	})

	return results
}

func (c Client) ok(accounts []Account) error {
	errCount := 0

	for _, acct := range accounts {
		if !acct.ok() {
			errCount++
		}
	}

	if errCount > 0 {
		return NewError(nil, "%d error(s) / %d account(s)", errCount, len(accounts))
	}

	return nil
}

func (c Client) worker(cache DBCache, in chan *Account, out chan<- *Account) {
	for acct := range in {
		if ca, err := cache.get(c.params.CfgName, acct.Object); err == nil {
			acct.Error = nil
			acct.StatusCode = ca.StatusCode
			acct.Timestamp = ca.Timestamp
			acct.Value = ca.Value

			out <- acct

			continue
		}

		acct.newTry()

		c.get(acct)

		if !acct.ok() {
			c.params.Errorf("Failed to get %v", acct)

			if acct.retry(c.params.MaxTries) {
				go func(acct *Account) {
					time.Sleep(time.Duration(acct.Try*acct.Try) * c.params.Wait)
					in <- acct
				}(acct)

				continue
			}
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

	result.Set("AppID", c.params.AppID)
	result.Set("Safe", c.params.Safe)
	result.Set("Object", object)

	return result
}

func (c Client) get(acct *Account) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		c.url(c.query(acct.Object)).String(),
		nil,
	)
	if err != nil {
		acct.Error = err

		return
	}

	response, err := c.http.Do(req)
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
