package feed

import (
	"encoding/xml"
	"html/template"
)

type rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`
	// Optional
	PubDate  string `xml:"channel>pubDate"`
	ItemList []item `xml:"channel>item"`
}

type item struct {
	// Required
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
	// Optional
	Content    template.HTML `xml:"encoded"`
	PubDate    string        `xml:"pubDate"`
	Comments   string        `xml:"comments"`
	Categories []string      `xml:"category"`
}

type atom1 struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title     string   `xml:"title"`
	Subtitle  string   `xml:"subtitle"`
	ID        string   `xml:"id"`
	Updated   string   `xml:"updated"`
	Rights    string   `xml:"rights"`
	Link      link     `xml:"link"`
	Author    author   `xml:"author"`
	EntryList []entry  `xml:"entry"`
}

type link struct {
	Href string `xml:"href,attr"`
}

type author struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type entry struct {
	Title      string   `xml:"title"`
	Summary    string   `xml:"summary"`
	Content    string   `xml:"content"`
	ID         string   `xml:"id"`
	Updated    string   `xml:"updated"`
	Link       link     `xml:"link"`
	Author     author   `xml:"author"`
	Categories []string `xml:"category"`
}

func atom1ToRss2(a atom1) rss2 {
	r := rss2{
		Title:       a.Title,
		Link:        a.Link.Href,
		Description: a.Subtitle,
		PubDate:     a.Updated,
	}
	r.ItemList = make([]item, len(a.EntryList))
	for i, entry := range a.EntryList {
		r.ItemList[i].Title = entry.Title
		r.ItemList[i].Link = entry.Link.Href
		if entry.Content == "" {
			r.ItemList[i].Description = template.HTML(entry.Summary)
		} else {
			r.ItemList[i].Description = template.HTML(entry.Content)
		}
	}
	return r
}

const atomErrStr = "expected element type <rss> but have <feed>"

func parseAtom(content []byte) (rss2, bool) {
	a := atom1{}
	err := xml.Unmarshal(content, &a)
	if err != nil {
		return rss2{}, false
	}
	return atom1ToRss2(a), true
}

func parseFeedContent(content []byte) (rss2, bool) {
	v := rss2{}
	err := xml.Unmarshal(content, &v)
	if err != nil {
		if err.Error() == atomErrStr {
			// try Atom 1.0
			return parseAtom(content)
		}
		return v, false
	}

	if v.Version == "2.0" {
		// RSS 2.0
		for i := range v.ItemList {
			if v.ItemList[i].Content != "" {
				v.ItemList[i].Description = v.ItemList[i].Content
			}
		}
		return v, true
	}

	return v, false
}
