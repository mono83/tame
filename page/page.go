package page

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"time"
)

// Page represents HTTP response page
type Page struct {
	URL        *url.URL
	StatusCode int
	Header     http.Header
	Elapsed    time.Duration
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
		"HTML [%d] %s %s with %d bytes fetched in %.2fs",
		p.StatusCode,
		p.URL.String(),
		p.ContentType(),
		len(p.Body),
		p.Elapsed.Seconds(),
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
