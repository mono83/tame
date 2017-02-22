package user

import (
	"github.com/mono83/tame/page"
	"net/http"
	"net/http/httptrace"
	"time"
)

func addTracing(req *http.Request) (*http.Request, *page.Trace) {
	t := page.NewTrace()

	ht := &httptrace.ClientTrace{
		DNSDone: func(httptrace.DNSDoneInfo) {
			t.DNS = time.Now().Sub(t.Init)
		},
		ConnectDone: func(network, addr string, err error) {
			t.Connection = time.Now().Sub(t.Init)
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			t.RequestSent = time.Now().Sub(t.Init)
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), ht)), t
}
