package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mono83/tame"
	"github.com/mono83/tame/goquery"
)

var headHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Title of the document</title>
</head>

<body>
</body>

</html> 
`

func TestUnmarshalHead(t *testing.T) {
	doc, err := goquery.FromDocument(tame.ByteDocument([]byte(headHTML)))
	if assert.NoError(t, err) {
		var head Head
		if assert.NoError(t, tame.Unmarshal(doc, &head)) {
			assert.Equal(t, "Title of the document", head.Title)
		}
	}
}
