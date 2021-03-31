package ecb

import (
	"github.com/mono83/tame"
	"github.com/mono83/tame/client"
)

// Fetch reads data from remote server
// transport - transport to use. Optional, if nil given will utilize client.New()
// url       - remote URL. Optional, if empty will use default URL
func Fetch(transport interface {
	Get(addr string) (tame.Document, error)
}, url string) (*Rates, error) {
	r := Rates{}
	if transport == nil {
		transport = client.New()
	}
	if len(url) == 0 {
		url = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	}

	doc, err := transport.Get(url)
	if err != nil {
		return nil, err
	}
	if err := r.Unmarshal(doc.GetBody()); err != nil {
		return nil, err
	}

	return &r, nil
}
