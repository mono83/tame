package goquery

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mono83/tame"
)

// wrapper is a wrapper over PuerkitoBio goquery library
// for DOM manipulation
type wrapper struct {
	sel *goquery.Selection
}

func (w wrapper) Size() int {
	return w.sel.Size()
}

func (w wrapper) Find(query string) tame.DOMSelection {
	return wrapper{sel: w.sel.Find(query)}
}
func (w wrapper) Attr(attr, defaultValue string) string {
	return w.sel.AttrOr(attr, defaultValue)
}
func (w wrapper) ReadText(query string, target *string) error {
	found := w.sel.Find(query)
	if found.Size() != 1 {
		*target = ""
		return fmt.Errorf("expected 1 element but found %d", found.Size())
	}

	*target = strings.TrimSpace(found.Text())
	return nil
}
func (w wrapper) ReadAttr(query, attr string, target *string) error {
	found := w.sel.Find(query)
	if found.Size() != 1 {
		*target = ""
		return fmt.Errorf("expected 1 element but found %d", found.Size())
	}

	*target = strings.TrimSpace(found.AttrOr(attr, ""))
	return nil
}
func (w wrapper) Each(f func(t tame.DOMSelection)) {
	if f != nil {
		w.sel.Each(func(_ int, s *goquery.Selection) {
			f(wrapper{sel: s})
		})
	}
}
