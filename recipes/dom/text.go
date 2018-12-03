package dom

import (
	"errors"
	"github.com/mono83/tame"
	"regexp"
	"strings"
)

// Text is used to extract plain text from DOM element
type Text struct {
	Selector  string
	PlainText string
	Links     []struct {
		Text     string
		Location string
	}
}

// UnmarshalDOM fills struct contents from DOM source
func (t *Text) UnmarshalDOM(src tame.DOMSelection) error {
	if src == nil {
		return errors.New("empty DOM Selection")
	}
	if len(t.Selector) == 0 {
		return errors.New("empty selector to search")
	}

	x := src.Find(t.Selector)
	if x.Size() == 0 {
		return nil
	} else if x.Size() > 1 {
		return errors.New("more than one element found using selector " + t.Selector)
	}

	// Obtaining plain text
	t.PlainText = cleanupText(x.AsText())

	// Obtaining links
	x.Find("a").Each(func(a tame.DOMSelection) {
		txt := cleanupText(a.AsText())
		href := a.Attr("href", "")

		if len(txt) > 0 && len(href) > 0 {
			t.Links = append(t.Links, struct {
				Text     string
				Location string
			}{
				Text:     txt,
				Location: href,
			})
		}
	})

	return nil
}

var anyWhitespace = regexp.MustCompile(`[^\S\n]+`)
var anyQuotes = regexp.MustCompile(`['"«»“”]`)

func cleanupText(s string) string {
	s = strings.TrimSpace(s)
	s = anyWhitespace.ReplaceAllString(s, " ")
	s = anyQuotes.ReplaceAllString(s, `"`)
	return s
}
