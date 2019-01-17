package feed

import "github.com/mono83/tame/clean"

type genericXMLItem struct {
	Title    string   `xml:"title"`
	LinkHref linkHref `xml:"link"`
	Rating   string   `xml:"rating"`

	ShortSummary     string   `xml:"summary"`
	ShortDescription string   `xml:"description"`
	Content          string   `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Categories       []string `xml:"category"`

	Enclosures []genericXMLEnclosure `xml:"enclosure"`
}

func (g genericXMLItem) Link() string {
	return any(g.LinkHref.Value, g.LinkHref.Text)
}

func (g genericXMLItem) Summary() string {
	return any(g.ShortSummary, g.ShortDescription)
}

func (g genericXMLItem) ToItem() Item {
	i := Item{
		Title:   clean.Trim(g.Title),
		Link:    clean.Trim(g.Link()),
		Short:   clean.Trim(g.ShortSummary),
		Content: clean.Trim(g.ShortDescription),
		Tags:    clean.Strings(g.Categories, clean.Trim),
	}

	for _, c := range g.Enclosures {
		i.Enclosure = append(i.Enclosure, Enclosure{
			Link: clean.Trim(c.Link),
			Type: clean.Trim(c.Type),
		})
	}

	return i
}

type genericXMLEnclosure struct {
	Link string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}