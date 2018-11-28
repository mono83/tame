package dom

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/mono83/tame"
)

// Head contains information, extracted from HTML head
type Head struct {
	Title        string
	Description  string
	Engine       string
	Keywords     []string
	URLCanonical *url.URL
}

// KeywordsCS returns comma separated keywords
func (h *Head) KeywordsCS() string {
	return strings.Join(h.Keywords, ",")
}

// UnmarshalDOM fills struct contents from DOM source
func (h *Head) UnmarshalDOM(src tame.DOMSelection) error {
	if src == nil {
		return errors.New("empty DOM Selection")
	}

	// Reading head
	domHead := src.Find("head")
	if domHead.Size() != 1 {
		return fmt.Errorf("expected 1 <head> tag, but found %d", domHead.Size())
	}

	head := Head{}
	domHead.ReadText("title", &head.Title)
	domHead.ReadAttr("meta[name=\"generator\"]", "content", &head.Engine)
	domHead.ReadAttr("meta[name=\"description\"]", "content", &head.Description)

	// Base keywords
	var keywordsString string
	domHead.ReadAttr("meta[name=\"keywords\"]", "content", &keywordsString)
	if len(keywordsString) > 0 {
		var kw []string
		if strings.Index(keywordsString, ",") > 0 {
			kw = strings.Split(keywordsString, ",")
		} else {
			kw = strings.Split(keywordsString, " ")
		}

		for _, k := range kw {
			k = strings.TrimSpace(k)
			if len(k) > 0 {
				head.Keywords = append(head.Keywords, k)
			}
		}
	}
	// Article tags
	domHead.Find("meta[property=\"article:tag\"]").Each(func(s tame.DOMSelection) {
		if s := s.Attr("content", ""); s != "" {
			head.Keywords = append(head.Keywords, s)
		}
	})

	var canonical string
	domHead.ReadAttr("link[rel=\"canonical\"]", "href", &canonical)
	if len(canonical) > 0 {
		head.URLCanonical, _ = url.Parse(canonical)
	}

	*h = head
	return nil
}
