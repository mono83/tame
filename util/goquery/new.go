package goquery

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

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

// FromDocumentE builds goquery implementation of tame.DOMDocument interface
func FromDocumentE(d tame.Document, err error) (tame.DOMDocument, error) {
	if err != nil {
		return nil, err
	}

	return FromDocument(d)
}

// FromResponse builds goquery implementation of tame.DOMDocument interface
// from http.Response
func FromResponse(res *http.Response) (tame.DOMDocument, error) {
	if res == nil {
		return nil, errors.New("nil response")
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(bts))
	if err != nil {
		return nil, err
	}

	var url url.URL
	if res.Request != nil && res.Request.URL != nil {
		url = *res.Request.URL
	}

	return document{
		url:     url,
		headers: res.Header,
		body:    bts,
		wrapper: wrapper{sel: doc.Selection},
	}, nil
}
