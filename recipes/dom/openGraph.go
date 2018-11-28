package dom

import (
	"errors"
	"net/url"

	"github.com/mono83/tame"
)

// OpenGraph contains OpenGraph data
type OpenGraph struct {
	Type        string
	SiteName    string
	Title       string
	Description string
	Locale      string
	Images      []url.URL
	URL         *url.URL
}

// UnmarshalDOM fills data from DOM structure
func (o *OpenGraph) UnmarshalDOM(src tame.DOMSelection) error {
	if src == nil {
		return errors.New("empty DOM Selection")
	}

	head := src.Find("head")

	head.ReadAttr("meta[property=\"og:type\"]", "content", &o.Type)
	head.ReadAttr("meta[property=\"og:title\"]", "content", &o.Title)
	head.ReadAttr("meta[property=\"og:locale\"]", "content", &o.Locale)
	head.ReadAttr("meta[property=\"og:site_name\"]", "content", &o.SiteName)
	head.ReadAttr("meta[property=\"og:description\"]", "content", &o.Description)

	var u string
	head.ReadAttr("meta[property=\"og:url\"]", "content", &u)
	if len(u) > 0 {
		o.URL, _ = url.Parse(u)
	}

	head.Find("meta[property=\"og:image\"]").Each(func(s tame.DOMSelection) {
		u := s.Attr("content", "")
		if len(u) > 0 {
			pu, err := url.Parse(u)
			if err == nil && pu != nil {
				o.Images = append(o.Images, *pu)
			}
		}
	})

	return nil
}
