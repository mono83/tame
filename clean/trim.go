package clean

import "strings"

// Trim strips leading and trailing spaces from string
func Trim(str string) string {
	return strings.TrimSpace(str)
}
