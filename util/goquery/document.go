package goquery

import (
	"net/http"
	"net/url"
)

// document is tame.DOMDocument interface implementation
type document struct {
	url     url.URL
	headers http.Header
	body    []byte
	wrapper
}

func (d document) GetURL() url.URL         { return d.url }
func (d document) GetHeaders() http.Header { return d.headers }
func (d document) GetBody() []byte         { return d.body }
