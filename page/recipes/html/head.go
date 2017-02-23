package html

import (
	"github.com/mono83/tame/page"
	"github.com/mono83/tame/page/recipes"
	"net/url"
	"strings"
)

var headRecipe = "HTMLHead"

// Head contains information, extracted from HTML head
type Head struct {
	Title        string
	Description  string
	Engine       string
	Keywords     []string
	URLCanonical *url.URL
	OpenGraph    OpenGraph
}

// OpenGraph contains OpenGraph data
type OpenGraph struct {
	Type        string
	Title       string
	Description string
	Images      []*url.URL
	URL         *url.URL
}

// HeadRecipe parses HTML head into html.Head struct
func HeadRecipe(p *page.Page, target interface{}) error {
	if p == nil {
		return recipes.Error(headRecipe, "Empty page provided")
	}
	ref, ok := target.(*Head)
	if !ok {
		return recipes.Error(headRecipe, "Recipe works only with *feed.Feed")
	}

	// Reading document
	doc, err := p.AsDOMSelection()
	if err != nil {
		return recipes.ErrorCaused(headRecipe, err)
	}

	// Reading head
	domHead := doc.Find("head")
	if domHead.Size() != 1 {
		return recipes.Errorf(headRecipe, "Expected 1 <head> tag, but found %d", domHead.Size())
	}

	head := Head{OpenGraph: OpenGraph{}}
	domHead.ReadText("title", &head.Title)
	domHead.ReadAttr("meta[name=\"generator\"]", "content", &head.Engine)
	domHead.ReadAttr("meta[name=\"description\"]", "content", &head.Description)

	// Base keywords
	var keywordsString string
	domHead.ReadAttr("meta[name=\"keywords\"]", "content", &keywordsString)
	if len(keywordsString) > 0 {
		head.Keywords = strings.Split(keywordsString, " ")
	}
	// Article tags
	domHead.Find("meta[property=\"article:tag\"]").Each(func(i int, s *page.DOMSelection) {
		if s, ok := s.Attr("content"); ok && s != "" {
			head.Keywords = append(head.Keywords, s)
		}
	})
	// Deduplication
	head.Keywords = removeDuplicatesUnordered(head.Keywords)

	var canonical string
	domHead.ReadAttr("link[rel=\"canonical\"]", "href", &canonical)
	if len(canonical) > 0 {
		head.URLCanonical, _ = url.Parse(canonical)
	}

	readOpenGraph(domHead, &head.OpenGraph)

	*ref = head
	return nil
}

func readOpenGraph(head *page.DOMSelection, og *OpenGraph) {
	head.ReadAttr("meta[property=\"og:type\"]", "content", &og.Type)
	head.ReadAttr("meta[property=\"og:title\"]", "content", &og.Title)
	head.ReadAttr("meta[property=\"og:description\"]", "content", &og.Description)

	var u string
	head.ReadAttr("meta[property=\"og:url\"]", "content", &u)
	if len(u) > 0 {
		og.URL, _ = url.Parse(u)
	}

	head.Find("meta[property=\"og:image\"]").Each(func(i int, s *page.DOMSelection) {
		u := s.AttrOr("content", "")
		if len(u) > 0 {
			pu, err := url.Parse(u)
			if err == nil && pu != nil {
				og.Images = append(og.Images, pu)
			}
		}
	})
}

func removeDuplicatesUnordered(elements []string) []string {
	if len(elements) == 0 {
		return []string{}
	}

	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
