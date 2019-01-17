package client

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"github.com/dsnet/compress/brotli"
	"github.com/mono83/tame/log"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mono83/tame"
	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
)

// Client represents new HTTP user session
type Client struct {
	// User-Agent, used by this client
	UserAgent string
	// HTTP Referer
	Referer string
	// Ray is xray.Ray logger
	Ray xray.Ray
	// Cookies
	cookies map[string][]*http.Cookie
	// Other HTTP headers
	Header map[string]string

	m      sync.Mutex
	client *http.Client
}

// New creates new HTTP user.
func New() *Client {
	c := &Client{
		Header: map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Encoding": "gzip, deflate, sdch, br",
			"Accept-Language": "en-US,en;q=0.8,ru;q=0.6,uk;q=0.4,pl;q=0.2",
		},
		Ray:       xray.ROOT.Fork().WithLogger("tame.client"),
		cookies:   map[string][]*http.Cookie{},
		client:    &http.Client{},
		UserAgent: CommonUserAgents[rand.Intn(len(CommonUserAgents))],
	}

	c.client.Jar = c
	return c
}

// NewRequest builds and returns new request of desired type with headers injected
func (c *Client) NewRequest(method, addr string, body io.Reader) (*http.Request, error) {
	ray := c.Ray.Fork().With(log.URL(addr), log.Method(method))

	ray.Debug("Building request :method :url")
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, err
	}

	// Injecting common headers
	for name, value := range c.Header {
		req.Header.Set(name, value)
		ray.Trace("Setting header :name = :value", args.Name(name), args.String{N: "value", V: value})
	}

	// Injecting user agent
	ray.Trace("Injecting user agent :name", args.Name(c.UserAgent))
	req.Header.Set("User-Agent", c.UserAgent)

	// Accepting gzip
	ray.Trace("Injecting accept encoding")
	req.Header.Add("Accept-Encoding", "gzip")

	// Injecting referer
	if len(c.Referer) > 0 {
		ray.Debug("Setting referer :referer", args.String{N: "referer", V: c.Referer})
		req.Header.Set("Referer", c.Referer)
	}

	return req, nil
}

// Get performs GET request
func (c *Client) Get(addr string) (tame.Document, error) {
	if len(addr) == 0 {
		return nil, errors.New("empty remote address")
	}
	ray := c.Ray.Fork().With(log.URL(addr))

	// Building request
	req, err := c.NewRequest("GET", addr, nil)
	if err != nil {
		ray.Error("Error building GET request :url - :err", args.Error{Err: err})
		return nil, err
	}

	// Sending request
	before := time.Now()
	ray.Debug("Sending request to :addr")
	resp, err := c.client.Do(req)
	if err != nil {
		ray.Error("Error performing GET :url - :err", args.Error{Err: err})
		return nil, err
	}
	defer silent(resp.Body.Close())

	// Checking against compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		ray.Trace("Detected GZIP encoding")
		reader, err = gzip.NewReader(resp.Body)
		defer silent(reader.Close())
	case "deflate":
		ray.Trace("Detected DEFLATE encoding")
		reader, err = zlib.NewReader(resp.Body)
		defer silent(reader.Close())
	case "br":
		ray.Trace("Detected BROTLI encoding")
		reader, err = brotli.NewReader(resp.Body, nil)
		defer silent(reader.Close())
	default:
		reader = resp.Body
	}
	if err != nil {
		ray.Error("Unable to establish reader - :err", args.Error{Err: err})
		return nil, err
	}

	// Reading response
	doc := document{
		url:    *req.URL,
		header: resp.Header,
		code:   resp.StatusCode,
	}

	doc.body, err = ioutil.ReadAll(reader)
	if err != nil {
		ray.Error("Unable to read response body for :url - :err", args.Error{Err: err})
		return nil, err
	}

	delta := time.Now().Sub(before)

	ray.Debug(
		"Document received with :code code, :count bytes in :delta",
		args.Int{N: "code", V: doc.code},
		args.Count(len(doc.body)),
		args.Delta(delta),
	)

	ray.InBytes(doc.body)

	// Replacing URL if working using own buffer proxy
	if xt := resp.Header.Get("x-tame-original"); xt != "" {
		xtu, xterr := url.Parse(xt)
		if xterr == nil {
			doc.url = *xtu
		}
	}

	return doc, nil
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (c *Client) SetCookies(url *url.URL, cookies []*http.Cookie) {
	c.m.Lock()
	c.cookies[url.Host] = cookies
	// Invalidating old cookies
	now := time.Now()
	for h, cs := range c.cookies {
		var newList []*http.Cookie
		for _, c := range cs {
			if c.Expires.After(now) || len(c.RawExpires) == 0 {
				newList = append(newList, c)
			}
		}
		c.cookies[h] = newList
	}
	c.m.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (c *Client) Cookies(url *url.URL) []*http.Cookie {
	return c.cookies[url.Host]
}

// GetCookie returns cookie by its name
func (c *Client) GetCookie(url *url.URL, name string) (string, bool) {
	if url == nil || len(name) == 0 {
		return "", false
	}

	c.m.Lock()
	defer c.m.Unlock()

	list, ok := c.cookies[url.Host]
	if !ok {
		return "", false
	}

	name = strings.ToLower(name)
	for _, c := range list {
		if name == strings.ToLower(c.Name) {
			return c.Value, true
		}
	}

	return "", false
}

func silent(error) {
}
