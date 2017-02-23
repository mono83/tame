package page

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// AsDOMSelection returns page contents as DOM tree, parsed using goquery.
// This function parses document on first demand and caches it.
func (p Page) AsDOMSelection() (*DOMSelection, error) {
	if p.document == nil {
		var err error
		p.document, err = goquery.NewDocumentFromReader(bytes.NewBuffer(p.Body))
		if err != nil {
			return nil, err
		}
	}

	return &DOMSelection{p.document.Selection}, nil
}

// DOMSelection is abstraction over goquery.Selection
type DOMSelection struct {
	*goquery.Selection
}

// Find gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new DOMSelection object
// containing these matched elements.
func (d *DOMSelection) Find(query string) *DOMSelection {
	return &DOMSelection{d.Selection.Find(query)}
}

// Each iterates over a Selection object, executing a function for each
// matched element. It returns the current DOMSelection object. The function
// f is called for each element in the selection with the index of the
// element in that selection starting at 0, and a *DOMSelection that contains
// only that element.
func (d *DOMSelection) Each(f func(int, *DOMSelection)) *DOMSelection {
	d.Selection.Each(func(i int, s *goquery.Selection) {
		f(i, &DOMSelection{s})
	})

	return d
}

// ReadText reads text from node found by query and writes it into target
func (d *DOMSelection) ReadText(query string, target *string) error {
	found := d.Selection.Find(query)
	if found.Size() != 1 {
		*target = ""
		return fmt.Errorf("Expected 1 element but found %d", found.Size())
	}

	*target = strings.TrimSpace(found.Text())
	return nil
}

// ReadAttr reads required attribute from node found by query and writes it into target
func (d *DOMSelection) ReadAttr(query, attr string, target *string) error {
	found := d.Selection.Find(query)
	if found.Size() != 1 {
		*target = ""
		return fmt.Errorf("Expected 1 element but found %d", found.Size())
	}

	*target = strings.TrimSpace(found.AttrOr(attr, ""))
	return nil
}
