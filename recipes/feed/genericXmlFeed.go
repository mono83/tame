package feed

import (
	"encoding/xml"
	"github.com/mono83/tame/clean"
)

type genericXMLFeed struct {
	TitleChannel string `xml:"channel>title"`
	TitleFeed    string `xml:"title"`

	LinkChannel string   `xml:"channel>link"`
	LinkFeed    linkHref `xml:"link"`

	Description string `xml:"channel>description"`
	Language    string `xml:"channel>language"`
	Updated     string `xml:"updated"`

	ItemsRSS  []genericXMLItem `xml:"channel>item"`
	ItemsAtom []genericXMLItem `xml:"entry"`
}

func (g genericXMLFeed) ToFeed() Feed {
	f := Feed{
		Title:       clean.Trim(g.Title()),
		Link:        clean.Trim(g.Link()),
		Description: clean.Trim(g.Description),
		Language:    clean.Trim(g.Language),
	}

	for _, i := range g.Items() {
		f.Items = append(f.Items, i.ToItem())
	}

	return f
}

func (g genericXMLFeed) Title() string {
	return any(g.TitleChannel, g.TitleFeed)
}

func (g genericXMLFeed) Link() string {
	return any(g.LinkChannel, g.LinkFeed.Value)
}

func (g genericXMLFeed) Items() []genericXMLItem {
	if len(g.ItemsRSS) > 0 {
		return g.ItemsRSS
	}

	return g.ItemsAtom
}

type linkHref struct {
	Text  string `xml:",chardata"`
	Value string `xml:"href,attr"`
}

func parseGeneric(src []byte) (*genericXMLFeed, error) {
	var v genericXMLFeed
	err := xml.Unmarshal(src, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func any(a, b string) string {
	if len(a) > 0 {
		return a
	}

	return b
}
