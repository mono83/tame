package clean

import (
	"net/url"
	"strings"
)

// UTMMarks removes UTM markers from link (if found)
func UTMMarks(original string) (result string) {
	result = strings.TrimSpace(original)
	if len(result) > 0 {
		if u, e := url.Parse(result); e == nil {
			if q, e := url.ParseQuery(u.RawQuery); e == nil {
				q.Del("utm_source")
				q.Del("utm_medium")
				q.Del("utm_campaign")
				q.Del("utm_term")
				q.Del("utm_content")

				u.RawQuery = q.Encode()
				result = u.String()
			}
		}
	}
	return
}
