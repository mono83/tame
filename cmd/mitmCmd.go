package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mono83/xray/args"

	"github.com/mono83/xray"

	"github.com/spf13/cobra"
)

var mitmCmd = &cobra.Command{
	Use:   "mitm port [folder]",
	Short: "Establishes mitm proxy for requests like http://localhost:<port>/http://google.com",
	RunE: func(cmd *cobra.Command, a []string) error {
		static := ""
		if len(a) < 1 {
			cmd.Usage()
			return errors.New("Not enough arguments")
		} else if len(a) == 2 {
			static = a[1]
			if !strings.HasSuffix(static, "/") {
				static += "/"
			}
		}

		addr := a[0]

		if len(static) > 0 {
			xray.BOOT.Info("Handing static files from :folder", args.String{N: "folder", V: static})
		}

		xray.BOOT.Info("Running listening MITM server on :addr", args.String{N: "addr", V: addr})

		// Running web server
		return http.ListenAndServe(":"+a[0], mitmHandler{Folder: static})
	},
}

type mitmHandler struct {
	Folder string
}

func (mh mitmHandler) serveStatic(uri string, w http.ResponseWriter, r *http.Request) {
	log := xray.ROOT.Fork().WithLogger("tame-mitm").With(args.String{N: "uri", V: uri})

	// Reading file
	filename := mh.Folder + uri
	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		bts, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Error(
				"Unable to read static file :uri - :err",
				args.Error{Err: err},
			)
			httpWriteError(w, err)
			return
		}

		// Choosing content type
		if strings.HasSuffix(uri, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(uri, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}

		w.Write(bts)
	} else {
		log.Error(
			"Unable to locate static file :uri",
		)
		w.WriteHeader(404)
	}
	return
}

func (mh mitmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI[1:]
	log := xray.ROOT.Fork().WithLogger("tame-mitm").With(args.String{N: "uri", V: uri})
	log.Debug("Incoming request :uri")

	if "/favicon.ico" == r.RequestURI || "/robots.txt" == r.RequestURI || "/service-worker.js" == r.RequestURI {
		w.WriteHeader(404)
		log.Debug("Restricted URL")
		return
	}

	// Static files serving
	if !strings.HasPrefix(uri, "http") && len(mh.Folder) > 0 {
		mh.serveStatic(uri, w, r)
		return
	}

	// Making request
	req, err := http.NewRequest(r.Method, uri, nil)
	if err != nil {
		log.Error("Unable to build request - :err", args.Error{Err: err})
		httpWriteError(w, err)
		return
	}
	for h, vv := range r.Header {
		lh := strings.ToLower(h)
		if strings.HasPrefix(lh, "mitm-") {
			for _, v := range vv {
				log.Debug(
					"Adding header :name = :value",
					args.String{N: "name", V: h[5:]},
					args.String{N: "value", V: v},
				)
				req.Header.Add(h[5:], v)
			}
		}
	}
	before := time.Now()
	resp, err := http.DefaultClient.Do(req)

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

	log.Debug(
		"Page :uri obtained with :count bytes in :delta",
		args.Count(len(bts)),
		args.Delta(after.Nanoseconds()),
	)

	// Page already loaded
	w.WriteHeader(resp.StatusCode)

	// Writing headers
	for h, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(h, v)
		}
	}

	// Writing body
	w.Write(bts)
	return
}
