package tame

// DOMSelection describes DOM node
type DOMSelection interface {
	Size() int
	Find(query string) DOMSelection
	Attr(attr, defaultValue string) string

	Each(func(DOMSelection))

	ReadText(query string, target *string) error
	ReadAttr(query, attr string, target *string) error

	AsText() string
}

// DOMDocument is an extended version of document, that
// supports DOM operations
type DOMDocument interface {
	Document
	DOMSelection
}
