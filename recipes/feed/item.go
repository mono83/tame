package feed

import (
	"time"

	"github.com/mono83/tame/util/clean"
)

// Item contains feed item content.
type Item struct {
	Title     string      // Feed item title
	Link      string      // Feed item link
	Short     string      // Feed item short content
	Content   string      // Feed item full content
	Tags      []string    // Feed item categories
	Enclosure []Enclosure // Feed item enclosures

	PublishedAt time.Time // Feed item publish time. May be zero, if publish time wasn't read
}

// Enclosure represents media enclosures
type Enclosure struct {
	Link string // Media location
	Type string // Media type
}

// CleanUTM returns new Item instance with UTM markers cleared from link
func (i Item) CleanUTM() Item {
	return Item{
		Title:       i.Title,
		Link:        clean.UTMMarks(i.Link),
		Short:       i.Short,
		Content:     i.Content,
		Tags:        i.Tags,
		Enclosure:   i.Enclosure,
		PublishedAt: i.PublishedAt,
	}
}

// GetPublishedAtOrNow returns item publishing time or current timestamp
// if publish time is empty
func (i Item) GetPublishedAtOrNow() time.Time {
	if i.PublishedAt.IsZero() {
		return time.Now()
	}

	return i.PublishedAt
}
