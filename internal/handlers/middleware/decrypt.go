package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/fishus/go-advanced-metrics/internal/cryptokey"
)

func Decrypt(privateKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if len(privateKey) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")

			body, err := io.ReadAll(r.Body)
			if err != nil {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, err.Error(), http.StatusBadRequest)
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				return
			}

			body, err = cryptokey.Decrypt(body, privateKey)
			if err != nil {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, "Failed to decrypt request data", http.StatusBadRequest)
				} else {
					http.Error(w, "Failed to decrypt request data", http.StatusBadRequest)
				}
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
