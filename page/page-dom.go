package page

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

// AsDOMDocument returns page contents as DOM tree, parsed using goquery.
// This function parses document on first demand and caches it.
func (p Page) AsDOMDocument() (*goquery.Document, error) {
	if p.document == nil {
		var err error
		p.document, err = goquery.NewDocumentFromReader(bytes.NewBuffer(p.Body))
		if err != nil {
			return nil, err
		}
	}

	return p.document, nil
}
