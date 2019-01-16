package feed

import "github.com/mono83/tame/clean"

// Item contains feed item content.
type Item struct {
	Title   string
	Link    string
	Content string
	Tags    []string
}

// CleanUTM returns new Item instance with UTM markers cleared from link
func (i Item) CleanUTM() Item {
	return Item{
		Title:   i.Title,
		Link:    clean.UTMMarks(i.Link),
		Content: i.Content,
		Tags:    i.Tags,
	}
}
