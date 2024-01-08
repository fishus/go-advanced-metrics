package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Decompress(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
