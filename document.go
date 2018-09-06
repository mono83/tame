package tame

import (
	"net/http"
	"net/url"
)

// Document interface describes parseable entities (HTML pages in common)
type Document interface {
	GetURL() url.URL
	GetHeaders() http.Header
	GetBody() []byte
}

// ByteDocument is simple implementation on Document interface
// Also it is plain wrapper over byte slice
type ByteDocument []byte

// GetURL returns empty URL
func (ByteDocument) GetURL() url.URL { return url.URL{} }

// GetHeaders returns empty headers map
func (ByteDocument) GetHeaders() http.Header { return http.Header{} }

// GetBody returns body contents
func (b ByteDocument) GetBody() []byte { return b }
