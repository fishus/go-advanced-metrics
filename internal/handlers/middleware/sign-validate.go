package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/secure"
)

func ValidateSign(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			hashString := r.Header.Get("HashSHA256")
			headerHash, err := hex.DecodeString(hashString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			bodyHash := secure.Hash(body, key)
			if !bytes.Equal(headerHash, bodyHash) {
				http.Error(w, "Data integrity has been compromised", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
