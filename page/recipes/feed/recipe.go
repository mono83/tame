package feed

import (
	"github.com/mono83/tame/page"
	"github.com/mono83/tame/page/recipes"
)

var recipeErrors = recipes.ErrorBuilder{RecipeName: "feed"}

// Recipe is function, used to parse RSS content from page.
func Recipe(p *page.Page, target interface{}) error {
	if p == nil {
		return recipeErrors.ErrorEmptyPage()
	}
	ref, ok := target.(*Feed)
	if !ok {
		return recipeErrors.Error("Recipe works only with *feed.Feed")
	}

	// Parsing ATOM contents
	r, ok := parseAtom(p.Body)
	if !ok {
		// Parsing RSS2 contents
		r, ok = parseFeedContent(p.Body)
		if !ok {
			return recipeErrors.Error("Unable to parse feed content")
		}
	}

	f := Feed{
		Title:       r.Title,
		Link:        r.Link,
		Description: r.Description,
		Items:       []Item{},
	}

	for _, i := range r.ItemList {
		newItem := Item{
			Title:   i.Title,
			Link:    i.Link,
			Tags:    i.Categories,
			Content: string(i.Content),
		}

		if len(i.Content) == 0 {
			newItem.Content = string(i.Description)
		}

		f.Items = append(f.Items, newItem)
	}

	*ref = f
	return nil
}
