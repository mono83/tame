package clean

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var utmdata = []struct {
	Expected, Given string
}{
	{"http://example.com/", "http://example.com/"},
	{"http://example.com/", " http://example.com/ "},
	{"http://example.com/?id=1", " http://example.com/?id=1"},
	{"http://example.com/?id=1", " http://example.com/?id=1&utm_source=some"},
	{"http://example.com/?id=1", " http://example.com/?utm_medium=some&id=1"},
	{"http://example.com/?id=1", " http://example.com/?utm_campaign=some&id=1"},
	{"http://example.com/?id=1", " http://example.com/?utm_campaign=some&id=1&utm_term=tag"},
	{"http://example.com/?id=1", " http://example.com/?utm_content=some&id=1&utm_term=tag"},
}

func TestUTMMarks(t *testing.T) {
	for _, data := range utmdata {
		t.Run(data.Given, func(t *testing.T) {
			assert.Equal(t, data.Expected, UTMMarks(data.Given))
		})
	}
}
