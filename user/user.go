package user

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"github.com/mono83/slf/wd"
	"github.com/mono83/tame/page"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// User represents new HTTP user.
type User struct {
	// User-Agent, used by this user
	UserAgent string
	// HTTP Referer
	Referer string
	// Other HTTP headers
	Header map[string]string

	client http.Client
	log    wd.Watchdog
}

// New creates new HTTP user.
func New() *User {
	return &User{
		Header: map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Encoding": "gzip, deflate, sdch, br",
			"Accept-Language": "en-US,en;q=0.8,ru;q=0.6,uk;q=0.4,pl;q=0.2",
		},
		UserAgent: userAgents[rand.Intn(len(userAgents))],
		client:    http.Client{},
		log:       wd.NewLogger("user"),
	}
}

// NewRequest builds and returns new request of desired type with headers injected
func (u *User) NewRequest(method, addr string, body io.Reader) (*http.Request, error) {
	u.log.Debug("Building request :method :addr", wd.StringParam("method", method), wd.StringParam("addr", addr))
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, err
	}

	// Injecting common headers
	for name, value := range u.Header {
		req.Header.Set(name, value)
	}

	// Injecting user agent
	req.Header.Set("User-Agent", u.UserAgent)

	// Injecting referer
	if len(u.Referer) > 0 {
		req.Header.Set("Referer", u.Referer)
	}

	return req, nil
}

// Get performs GET request
func (u *User) Get(addr string) (*page.Page, error) {
	if len(addr) == 0 {
		return nil, errors.New("Empty remote address")
	}
	log := u.log.WithParams(wd.StringParam("addr", addr))

	// Building request
	req, err := u.NewRequest("GET", addr, nil)
	if err != nil {
		log.Error("Error building GET request :addr - :err", wd.ErrParam(err))
		return nil, err
	}

	// Sending request
	log.Debug("Sending request to :addr")
	before := time.Now()
	resp, err := u.client.Do(req)
	if err != nil {
		log.Error("Error while GET :addr - :err", wd.ErrParam(err))
		return nil, err
	}
	defer resp.Body.Close()

	// Checking against compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	case "deflate":
		reader, err = zlib.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	if err != nil {
		log.Error("Unable to establish reader - :err", wd.ErrParam(err))
		return nil, err
	}

	// Reading response
	p := &page.Page{
		Elapsed:    time.Now().Sub(before),
		URL:        req.URL,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}

	p.Body, err = ioutil.ReadAll(reader)
	if err != nil {
		log.Error("Unable to read response body for :addr - :err", wd.ErrParam(err))
		return nil, err
	}

	return p, nil
}
