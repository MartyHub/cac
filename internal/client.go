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
	cache  *Cache
	http   *http.Client
	log    *log.Logger // help testing
	params Parameters
}

func NewClient(params Parameters) (Client, error) {
	cert, err := tls.LoadX509KeyPair(params.CertFile, params.KeyFile)
	if err != nil {
		return Client{}, err
	}

	cache, err := NewCache(params.CfgName)
	if err != nil {
		return Client{}, err
	}

	clock := utcClock{}

	cache.clean(clock, params.Expiry)

	return Client{
		cache: cache,
		clock: clock,
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
	size := c.poolSize()
	in := make(chan *account, size)
	out := make(chan *account, size)

	for i := 0; i < size; i++ {
		go c.worker(in, out)
	}

	count := make(chan int)

	go c.read(in, count)

	accounts := c.collect(out, count)

	close(in)

	err := c.output(accounts)
	if err != nil {
		return err
	}

	if c.cache.exists() {
		c.cache.merge(accounts)

		return c.cache.save()
	}

	return c.ok(accounts)
}

func (c Client) output(accounts []account) error {
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

func (c Client) read(in chan<- *account, count chan<- int) {
	if c.params.fromStdin() {
		count <- c.readFromReader(in, os.Stdin)
	} else {
		count <- c.readFromParams(in)
	}
}

func (c Client) readFromParams(in chan<- *account) int {
	now := c.clock.now()
	result := 0

	for _, object := range c.params.Objects {
		in <- newAccount(object, now, "", "", "")
		result++
	}

	return result
}

func (c Client) readFromReader(in chan<- *account, reader io.Reader) int {
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

func (c Client) collect(accounts <-chan *account, count <-chan int) []account {
	l := 0
	lenKnown := false
	results := make([]account, 0)

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

func (c Client) ok(accounts []account) error {
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

func (c Client) worker(in chan *account, out chan<- *account) {
	for acct := range in {
		if ca, found := c.cache.Accounts[acct.Object]; found {
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
				go func() {
					time.Sleep(time.Duration(acct.Try*acct.Try) * c.params.Wait)
					in <- acct
				}()

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

func (c Client) get(acct *account) {
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
