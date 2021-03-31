package ecb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var srcXml = `
<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender>
		<gesmes:name>European Central Bank</gesmes:name>
	</gesmes:Sender>
	<Cube>
		<Cube time='2021-03-30'>
			<Cube currency='USD' rate='1.1741'/>
			<Cube currency='JPY' rate='129.48'/>
			<Cube currency='BGN' rate='1.9558'/>
			<Cube currency='CZK' rate='26.122'/>
			<Cube currency='DKK' rate='7.4369'/>
			<Cube currency='GBP' rate='0.85378'/>
			<Cube currency='HUF' rate='363.30'/>
			<Cube currency='PLN' rate='4.6582'/>
			<Cube currency='RON' rate='4.9210'/>
			<Cube currency='SEK' rate='10.2473'/>
			<Cube currency='CHF' rate='1.1057'/>
			<Cube currency='ISK' rate='148.50'/>
			<Cube currency='NOK' rate='10.0613'/>
			<Cube currency='HRK' rate='7.5698'/>
			<Cube currency='RUB' rate='89.1591'/>
			<Cube currency='TRY' rate='9.7800'/>
			<Cube currency='AUD' rate='1.5419'/>
			<Cube currency='BRL' rate='6.7685'/>
			<Cube currency='CAD' rate='1.4814'/>
			<Cube currency='CNY' rate='7.7154'/>
			<Cube currency='HKD' rate='9.1283'/>
			<Cube currency='IDR' rate='17063.20'/>
			<Cube currency='ILS' rate='3.9133'/>
			<Cube currency='INR' rate='86.2540'/>
			<Cube currency='KRW' rate='1331.35'/>
			<Cube currency='MXN' rate='24.2262'/>
			<Cube currency='MYR' rate='4.8737'/>
			<Cube currency='NZD' rate='1.6794'/>
			<Cube currency='PHP' rate='57.015'/>
			<Cube currency='SGD' rate='1.5815'/>
			<Cube currency='THB' rate='36.714'/>
			<Cube currency='ZAR' rate='17.5396'/>
		</Cube>
	</Cube>
</gesmes:Envelope>
`

func TestRates_Unmarshal(t *testing.T) {
	rates := Rates{}
	if assert.NoError(t, rates.Unmarshal([]byte(srcXml))) {
		assert.Equal(t, time.Date(2021, 3, 30, 0, 0, 0, 0, time.UTC), rates.Date)
		if assert.Len(t, rates.Rates, 32) {
			assert.Equal(t, "USD", rates.Rates[0].ISO)
			assert.Equal(t, float32(1.1741), rates.Rates[0].Rate)
		}
	}
}
