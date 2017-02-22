package recipes

import "github.com/mono83/tame/page"

// Recipe describes accessor function, used to parse document contents
type Recipe func(*page.Page, interface{}) error
