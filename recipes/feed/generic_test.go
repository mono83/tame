package feed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var genericRssSrc = `
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

</rss>`

var genericZenSrc = `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
    xmlns:content="http://purl.org/rss/1.0/modules/content/"
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:media="http://search.yahoo.com/mrss/"
    xmlns:atom="http://www.w3.org/2005/Atom"
    xmlns:georss="http://www.georss.org/georss">
    <channel>
        <title>Пастернак</title>
        <link>http://example.com/</link>
        <description>Проект о фруктах и овощах. Рассказываем о том, как выращивать, готовить и правильно есть.</description>
        <language>ru</language>
        <item>
           <title>Андроид восстановит ферму в Японии</title>
           <link>http://example.com/2023/07/04/android-happy-farmer</link>
           <pdalink>http://m.example.com/2023/07/04/android-happy-farmer</pdalink>
           <amplink>http://amp.example.com/2023/07/04/android-happy-farmer</amplink>
           <guid>2fd4e1c67a2d28fced849ee1bb76e7391b93eb12</guid>
           <pubDate>Tue, 4 Jul 2023 04:20:00 +0300</pubDate>
           <media:rating scheme="urn:simple">nonadult</media:rating>
           <author>Петр Стругацкий</author>
           <category>Технологии</category>
           <enclosure url="http://example.com/2023/07/04/pic1.jpg" type="image/jpeg"/>
           <enclosure url="http://example.com/2023/07/04/pic2.jpg" type="image/jpeg"/>
           <enclosure url="http://example.com/2023/07/04/video/420"
                      type="video/x-ms-asf"/>
           <description>
                <![CDATA[
Заброшенную землю рядом с токийским университетом Нисёгакуся передали андроиду
с внешностью известного японского хозяйственника.
]]>
            </description>
            <content:encoded>
                <![CDATA[

<p>Здесь находится полный текст статьи.
Этот текст может прерываться картинками, видео и другим медиа-контентом так же,
как в оригинальной статье. Пример вставленной картинки ниже.</p>
<figure>
    <img src="http://example.com/2023/07/04/pic1.jpg" width="1200" height="900">
        <figcaption>
Первый андроид-фермер смотрит на свои угодья

            <span class="copyright">Михаил Родченков</span>
        </figcaption>
    </figure>
    <p>Продолжение статьи после вставленной картинки. В статье рассказывается
о технологии вспахивании земли, которую использует японский андроид-фермер.
Поэтому в материале не обойтись без видеоролика. Пример видеоролика ниже.</p>
    <figure>
        <video width="1200" height="900">
            <source src="http://example.com/2023/07/04/video/42420" type="video/mp4">
            </video>
            <figcaption>
Андроид-фермер вспахивает землю при помощи собственного изобретения

                <span class="copyright">Михаил Родченков</span>
            </figcaption>
        </figure>
        <p>Статья продолжается после видео. Андроид копает картошку.
Фермы развиваются. Япония продолжает удивлять.</p>
]]>
            </content:encoded>
        </item>
    </channel>
</rss>
`

var genericAtomSrc = `
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

</feed>
`

func TestParseGenericRss(t *testing.T) {
	if f, err := parseGeneric([]byte(genericRssSrc)); assert.NoError(t, err) {
		assert.Equal(t, "W3Schools Home Page", f.Title())
		assert.Equal(t, "https://www.w3schools.com", f.Link())
		assert.Equal(t, "Free web building tutorials", f.Description)
		assert.Equal(t, "", f.Updated)
		assert.Equal(t, "", f.Language)

		if assert.Len(t, f.Items(), 2) {
			assert.Equal(t, "RSS Tutorial", f.Items()[0].Title)
			assert.Equal(t, "https://www.w3schools.com/xml/xml_rss.asp", f.Items()[0].Link())
			assert.Equal(t, "New RSS tutorial on W3Schools", f.Items()[0].Summary())

			assert.Equal(t, "XML Tutorial", f.Items()[1].Title)
			assert.Equal(t, "https://www.w3schools.com/xml", f.Items()[1].Link())
			assert.Equal(t, "New XML tutorial on W3Schools", f.Items()[1].Summary())

			assert.Len(t, f.Items()[0].Enclosures, 0)
		}
	}
}

func TestParseGenericZen(t *testing.T) {
	if f, err := parseGeneric([]byte(genericZenSrc)); assert.NoError(t, err) {
		assert.Equal(t, "Пастернак", f.Title())
		assert.Equal(t, "http://example.com/", f.Link())
		assert.Equal(t, "Проект о фруктах и овощах. Рассказываем о том, как выращивать, готовить и правильно есть.", f.Description)
		assert.Equal(t, "", f.Updated)
		assert.Equal(t, "ru", f.Language)

		if assert.Len(t, f.Items(), 1) {
			assert.Equal(t, "Андроид восстановит ферму в Японии", f.Items()[0].Title)
			assert.Equal(t, "http://example.com/2023/07/04/android-happy-farmer", f.Items()[0].Link())
			assert.Equal(t, "\n                \nЗаброшенную землю рядом с токийским университетом Нисёгакуся передали андроиду\nс внешностью известного японского хозяйственника.\n\n            ", f.Items()[0].Summary())
			assert.Equal(t, "nonadult", f.Items()[0].Rating)
			assert.Equal(t, 1740, len(f.Items()[0].Content))

			if assert.Len(t, f.Items()[0].Enclosures, 3) {
				assert.Equal(t, "http://example.com/2023/07/04/pic1.jpg", f.Items()[0].Enclosures[0].Link)
				assert.Equal(t, "http://example.com/2023/07/04/video/420", f.Items()[0].Enclosures[2].Link)

				assert.Equal(t, "image/jpeg", f.Items()[0].Enclosures[0].Type)
				assert.Equal(t, "video/x-ms-asf", f.Items()[0].Enclosures[2].Type)
			}
		}
	}
}

func TestParseGenericAtom(t *testing.T) {
	if f, err := parseGeneric([]byte(genericAtomSrc)); assert.NoError(t, err) {
		assert.Equal(t, "Example Feed", f.Title())
		assert.Equal(t, "http://example.org/", f.Link())
		assert.Equal(t, "", f.Description)
		assert.Equal(t, "2003-12-13T18:30:02Z", f.Updated)
		assert.Equal(t, "", f.Language)

		if assert.Len(t, f.Items(), 1) {
			assert.Equal(t, "Atom-Powered Robots Run Amok", f.Items()[0].Title)
			assert.Equal(t, "http://example.org/2003/12/13/atom03", f.Items()[0].Link())
			assert.Equal(t, "Some text.", f.Items()[0].Summary())

			assert.Len(t, f.Items()[0].Enclosures, 0)
		}
	}
}
