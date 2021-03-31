package ecb

import (
	"encoding/xml"
	"time"
)

// Rates contains rates data obtained from ECB
type Rates struct {
	Date  time.Time
	Rates []Rate
}

// Rate contains rate from EUR
type Rate struct {
	ISO  string
	Rate float32
}

// Unmarshal fills data to rates from given XML bytes
func (r *Rates) Unmarshal(src []byte) error {
	var intermediate xmlStruct
	if err := xml.Unmarshal(src, &intermediate); err != nil {
		return err
	}

	for _, x := range intermediate.Outer.Inner.Rates {
		r.Rates = append(r.Rates, Rate{
			ISO:  x.ISO,
			Rate: x.Rate,
		})
	}

	if d, err := time.Parse("2006-01-02", intermediate.Outer.Inner.Date); err == nil {
		r.Date = d.UTC()
	} else {
		return err
	}

	return nil
}

type xmlStruct struct {
	XMLName xml.Name `xml:"Envelope"`
	Outer   struct {
		Inner struct {
			Date  string `xml:"time,attr"`
			Rates []struct {
				ISO  string  `xml:"currency,attr"`
				Rate float32 `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}
