package client

import (
	"net/http"
	"net/url"
)

type document struct {
	code   int
	url    url.URL
	header http.Header
	body   []byte
}

func (d document) GetURL() url.URL         { return d.url }
func (d document) GetHeaders() http.Header { return d.header }
func (d document) GetBody() []byte         { return d.body }
