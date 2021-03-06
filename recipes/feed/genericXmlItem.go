package feed

import (
	"time"

	"github.com/mono83/tame/util/clean"
)

type genericXMLItem struct {
	Title        string   `xml:"title"`
	LinkHref     linkHref `xml:"link"`
	LinkOriginal string   `xml:"origLink"`
	Rating       string   `xml:"rating"`

	PublishTimeString string `xml:"pubDate"`

	ShortSummary     string   `xml:"summary"`
	ShortDescription string   `xml:"description"`
	Content          string   `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	ContentYandex    string   `xml:"http://news.yandex.ru full-text"`
	Categories       []string `xml:"category"`

	Enclosures []genericXMLEnclosure `xml:"enclosure"`
}

func (g genericXMLItem) Link() string {
	return any(g.LinkOriginal, any(g.LinkHref.Value, g.LinkHref.Text))
}

func (g genericXMLItem) Summary() string {
	return any(g.ShortSummary, g.ShortDescription)
}

func (g genericXMLItem) ToItem() Item {
	i := Item{
		Title:   clean.Trim(g.Title),
		Link:    clean.Trim(g.Link()),
		Short:   clean.Trim(g.Summary()),
		Content: clean.Trim(g.Content),
		Tags:    clean.Strings(g.Categories, clean.Trim),
	}

	if len(i.Content) == 0 {
		i.Content = i.Short
	}

	if len(i.Content) == 0 && len(g.ContentYandex) > 0 {
		i.Content = clean.Trim(g.ContentYandex)
	}

	// Parsing publishing time
	if t, err := time.Parse(time.RFC1123, g.PublishTimeString); err == nil {
		i.PublishedAt = t.UTC()
	} else if t, err = time.Parse(time.RFC1123Z, g.PublishTimeString); err == nil {
		i.PublishedAt = t.UTC()
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
