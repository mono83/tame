package feed

import (
	"testing"

	"github.com/mono83/tame"
	"github.com/stretchr/testify/assert"
)

var srcRss2 = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
  </item>
  <item>
    <title>XML Tutorial</title>
    <link>https://www.w3schools.com/xml</link>
    <description>New XML tutorial on W3Schools</description>
  </item>
</channel>

</rss> 
`

var srcAtom = `
<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">

  <title>Example Feed</title>
  <link href="http://example.org/"/>
  <updated>2003-12-13T18:30:02Z</updated>
  <author>
    <name>John Doe</name>
  </author>
  <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>

  <entry>
    <title>Atom-Powered Robots Run Amok</title>
    <link href="http://example.org/2003/12/13/atom03"/>
    <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
    <updated>2003-12-13T18:30:02Z</updated>
    <summary>Some text.</summary>
  </entry>

</feed>`

func TestFeedUnmarshalRSS2(t *testing.T) {
	var feed Feed
	if assert.NoError(t, tame.Unmarshal(tame.ByteDocument([]byte(srcRss2)), &feed)) {
		assert.Equal(t, "W3Schools Home Page", feed.Title)
		assert.Equal(t, "https://www.w3schools.com", feed.Link)
		assert.Equal(t, "Free web building tutorials", feed.Description)
		if assert.Len(t, feed.Items, 2) {
			assert.Equal(t, "RSS Tutorial", feed.Items[0].Title)
			assert.Equal(t, "https://www.w3schools.com/xml/xml_rss.asp", feed.Items[0].Link)
			assert.Equal(t, "New RSS tutorial on W3Schools", feed.Items[0].Content)

			assert.Equal(t, "XML Tutorial", feed.Items[1].Title)
			assert.Equal(t, "https://www.w3schools.com/xml", feed.Items[1].Link)
			assert.Equal(t, "New XML tutorial on W3Schools", feed.Items[1].Content)
		}
	}
}

func TestFeedUnmarshalAtom(t *testing.T) {
	var feed Feed
	if assert.NoError(t, tame.Unmarshal(tame.ByteDocument([]byte(srcAtom)), &feed)) {
		assert.Equal(t, "Example Feed", feed.Title)
		assert.Equal(t, "http://example.org/", feed.Link)
		assert.Equal(t, "", feed.Description)
		if assert.Len(t, feed.Items, 1) {
			assert.Equal(t, "Atom-Powered Robots Run Amok", feed.Items[0].Title)
			assert.Equal(t, "http://example.org/2003/12/13/atom03", feed.Items[0].Link)
			assert.Equal(t, "Some text.", feed.Items[0].Content)
		}
	}
}
