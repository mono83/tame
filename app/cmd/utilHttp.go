package cmd

import (
	"net/http"
)

// httpWriteError writes plaintext error response with 500 status code
func httpWriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(err.Error()))
	return
}
