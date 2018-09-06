package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mono83/xray"
	"github.com/mono83/xray/args"
	"github.com/spf13/cobra"
)

var pages = map[string]page{}

var pageCmd = &cobra.Command{
	Use:   "page port",
	Short: "Establishes proxy for requests like http://localhost:<port>/http://google.com",
	RunE: func(cmd *cobra.Command, a []string) error {
		if len(a) != 1 {
			cmd.Usage()
			return errors.New("Not enough arguments")
		}

		addr := a[0]
		xray.BOOT.Info("Running listening MITM page server on :addr", args.String{N: "addr", V: addr})

		// Running web server
		return http.ListenAndServe(":"+a[0], handler{})
	},
}

type page struct {
	code    int
	headers http.Header
	body    []byte
}

type handler struct{}

func (s handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI[1:]
	log := xray.ROOT.Fork().WithLogger("tame-mitm").With(args.String{N: "uri", V: uri})
	log.Debug("Incoming request :uri")

	if "/favicon.ico" == r.RequestURI || "/robots.txt" == r.RequestURI {
		w.WriteHeader(404)
		log.Debug("Restricted URL")
		return
	}

	p, ok := pages[uri]
	if !ok {
		// Making request
		before := time.Now()
		resp, err := http.Get(uri)
		if err != nil {
			log.Error("Received error - :err", args.Error{Err: err})
			httpWriteError(w, err)
			return
		}

		// Reading response
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		after := time.Now().Sub(before)
		if err != nil {
			log.Error("Unable to read response body - :err", args.Error{Err: err})
			httpWriteError(w, err)
			return
		}

		// Building and saving page
		p = page{
			code:    resp.StatusCode,
			headers: resp.Header,
			body:    bts,
		}
		pages[uri] = p
		log.Debug(
			"Page :uri obtained with :count bytes in :delta",
			args.Count(len(bts)),
			args.Delta(after.Nanoseconds()),
		)

	} else {
		log.Debug("Page :uri featched from cache")
	}

	// Writing headers
	w.Header().Add("x-tame-original", uri)
	for h, vv := range p.headers {
		for _, v := range vv {
			w.Header().Add(h, v)
		}
	}
	// Page already loaded
	w.WriteHeader(p.code)

	// Writing body
	w.Write(p.body)
	return
}
