package cmd

import (
	"errors"
	"github.com/mono83/slf/wd"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var pages = map[string]page{}

var pageCmd = &cobra.Command{
	Use:   "page port",
	Short: "Establishes proxy for requests like http://localhost:<port>/http://google.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Usage()
			return errors.New("Not enough arguments")
		}

		addr := args[0]
		httpLog.Info("Running listening server on :addr", wd.StringParam("addr", addr))

		// Running web server
		return http.ListenAndServe(":"+args[0], handler{})
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
	httpLog.Debug("Incoming request :uri", wd.StringParam("uri", uri))

	if "/favicon.ico" == r.RequestURI || "/robots.txt" == r.RequestURI {
		w.WriteHeader(404)
		httpLog.Debug("Restricted URL")
		return
	}

	p, ok := pages[uri]
	if !ok {
		// Making request
		resp, err := http.Get(uri)
		if err != nil {
			httpLog.Error("Received error - :err", wd.ErrParam(err))
			httpWriteError(w, err)
			return
		}

		// Reading response
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			httpLog.Error("Unable to read response body - :err", wd.ErrParam(err))
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
		httpLog.Debug("Page :uri obtained with :count bytes", wd.StringParam("uri", uri), wd.CountParam(len(p.body)))
	} else {
		httpLog.Debug("Page :uri featched from cache", wd.StringParam("uri", uri))
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
