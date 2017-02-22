package feed

// Feed contains feed contents.
type Feed struct {
	Title       string
	Link        string
	Description string

	Items []Item
}
