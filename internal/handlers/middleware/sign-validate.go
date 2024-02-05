package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/fishus/go-advanced-metrics/internal/secure"
)

func ValidateSign(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			hashString := r.Header.Get("HashSHA256")

			if len(key) == 0 || hashString == "" {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")

			headerHash, err := hex.DecodeString(hashString)
			if err != nil {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, err.Error(), http.StatusBadRequest)
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, err.Error(), http.StatusBadRequest)
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				return
			}
			// Restore the io.ReadCloser to it's original state
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			bodyHash := secure.Hash(body, key)

			if !bytes.Equal(headerHash, bodyHash) {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, "Data integrity has been compromised", http.StatusBadRequest)
				} else {
					http.Error(w, "Data integrity has been compromised", http.StatusBadRequest)
				}
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
