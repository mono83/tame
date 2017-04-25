package page

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// Page represents HTTP response page
type Page struct {
	Trace Trace

	URL        *url.URL
	StatusCode int
	Header     http.Header
	Body       []byte

	document *goquery.Document
}

// IsEmpty returns true if page byte content is empty
func (p Page) IsEmpty() bool {
	return len(p.Body) == 0
}

// ContentType returns content type of page and "application/octet-stream" as default
func (p Page) ContentType() string {
	ct := p.Header.Get("content-type")
	if len(ct) == 0 {
		return "application/octet-stream"
	}

	return ct
}

// String returns short string representation of page
func (p Page) String() string {
	return fmt.Sprintf(
		"HTML [%d] %s %s with %d bytes fetched in [%s]",
		p.StatusCode,
		p.URL.String(),
		p.ContentType(),
		len(p.Body),
		p.Trace.String(),
	)
}

// AsString converts body to string
func (p Page) AsString() string {
	return string(p.Body)
}

// AsJSON unmarshals data bytes into struct using JSON
func (p Page) AsJSON(target interface{}) error {
	return json.Unmarshal(p.Body, target)
}

// MarshalJSON is implementation of JSON parsing helpers
func (p Page) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["URL"] = p.URL.String()
	m["StatusCode"] = p.StatusCode
	m["Header"] = p.Header
	m["Body"] = string(p.Body)

	return json.Marshal(m)
}

// UnmarshalJSON is implementation of JSON parsing helpers
func (p *Page) UnmarshalJSON(data []byte) error {
	var m struct {
		URL        string              `json:"URL"`
		StatusCode int                 `json:"StatusCode"`
		Header     map[string][]string `json:"Header"`
		Body       string              `json:"body"`
	}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	url, err := url.Parse(m.URL)
	if err != nil {
		return err
	}

	p.URL = url
	p.Header = http.Header(m.Header)
	p.StatusCode = m.StatusCode
	p.Body = []byte(m.Body)

	return nil
}
