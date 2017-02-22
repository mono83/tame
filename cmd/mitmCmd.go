package cmd

import (
	"errors"
	"github.com/mono83/slf/wd"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var mitmCmd = &cobra.Command{
	Use:   "mitm port [folder]",
	Short: "Establishes mitm proxy for requests like http://localhost:<port>/http://google.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		static := ""
		if len(args) < 1 {
			cmd.Usage()
			return errors.New("Not enough arguments")
		} else if len(args) == 2 {
			static = args[1]
			if !strings.HasSuffix(static, "/") {
				static += "/"
			}
		}

		addr := args[0]

		if len(static) > 0 {
			httpLog.Info("Handing static files from :folder", wd.StringParam("folder", static))
		}

		httpLog.Info("Running listening MITM server on :addr", wd.StringParam("addr", addr))

		// Running web server
		return http.ListenAndServe(":"+args[0], mitmHandler{Folder: static})
	},
}

type mitmHandler struct {
	Folder string
}

func (mh mitmHandler) serveStatic(uri string, w http.ResponseWriter, r *http.Request) {
	// Reading file
	filename := mh.Folder + uri
	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		bts, err := ioutil.ReadFile(filename)
		if err != nil {
			httpLog.Error(
				"Unable to read static file :uri - :err",
				wd.StringParam("uri", uri),
				wd.ErrParam(err),
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
		httpLog.Error("Unable to locate static file :uri", wd.StringParam("uri", uri))
		w.WriteHeader(404)
	}
	return
}

func (mh mitmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI[1:]
	httpLog.Debug("Incoming request :uri", wd.StringParam("uri", uri))

	if "/favicon.ico" == r.RequestURI || "/robots.txt" == r.RequestURI || "/service-worker.js" == r.RequestURI {
		w.WriteHeader(404)
		httpLog.Debug("Restricted URL")
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
		httpLog.Error("Unable to build request - :err", wd.ErrParam(err))
		httpWriteError(w, err)
		return
	}
	for h, vv := range r.Header {
		lh := strings.ToLower(h)
		if strings.HasPrefix(lh, "mitm-") {
			for _, v := range vv {
				httpLog.Debug(
					"Adding header :name = :value",
					wd.StringParam("name", h[5:]),
					wd.StringParam("value", v),
				)
				req.Header.Add(h[5:], v)
			}
		}
	}
	resp, err := http.DefaultClient.Do(req)
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

	httpLog.Debug("Page :uri obtained with :count bytes", wd.StringParam("uri", uri), wd.CountParam(len(bts)))

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
