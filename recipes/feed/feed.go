package feed

import (
	"github.com/mono83/tame/util/clean"
)

// Feed contains feed contents.
type Feed struct {
	Title       string
	Link        string
	Description string
	Language    string

	Items []Item
}

// Unmarshal fills Feed structure from byte source
func (f *Feed) Unmarshal(src []byte) error {
	g, err := parseGeneric(src)
	if err != nil {
		return err
	}

	fc := g.ToFeed()
	*f = fc
	return nil
}

// CleanUTM returns new Feed instance with UTM markers cleared from all links
func (f Feed) CleanUTM() Feed {
	n := Feed{
		Title:       f.Title,
		Link:        clean.UTMMarks(f.Link),
		Description: f.Description,
		Language:    f.Language,
	}

	for _, i := range f.Items {
		n.Items = append(n.Items, i.CleanUTM())
	}

	return n
}
