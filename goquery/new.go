package goquery

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/mono83/tame"
)

// FromDocument builds goquery implementation of tame.DOMDocument interface
func FromDocument(d tame.Document) (tame.DOMDocument, error) {
	if d == nil {
		return nil, errors.New("nil document")
	}
	if x, ok := d.(tame.DOMDocument); ok {
		return x, nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(d.GetBody()))
	if err != nil {
		return nil, err
	}

	return document{
		url:     d.GetURL(),
		headers: d.GetHeaders(),
		body:    d.GetBody(),
		wrapper: wrapper{sel: doc.Selection},
	}, nil
}
