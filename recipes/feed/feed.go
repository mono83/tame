package feed

import "errors"

// Feed contains feed contents.
type Feed struct {
	Title       string
	Link        string
	Description string

	Items []Item
}

// Unmarshal fills Feed structure from byte source
func (f *Feed) Unmarshal(src []byte) error {
	// Parsing ATOM contents
	r, ok := parseAtom(src)
	if !ok {
		// Parsing RSS2 contents
		r, ok = parseFeedContent(src)
		if !ok {
			return errors.New("unable to parse feed contents")
		}
	}

	fc := Feed{
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
		if len(i.FeedBurnerOriginal) > 0 {
			newItem.Link = i.FeedBurnerOriginal
		}

		fc.Items = append(fc.Items, newItem)
	}

	*f = fc
	return nil
}
