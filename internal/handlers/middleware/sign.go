package middleware

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/secure"
)

func Sign(key []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if len(key) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			sign := secure.NewSign(key)
			hw := &signWriter{ResponseWriter: w, Sign: sign, Buf: &bytes.Buffer{}}
			defer hw.Close()

			next.ServeHTTP(hw, r)
		}
		return http.HandlerFunc(fn)
	}
}

type signWriter struct {
	http.ResponseWriter
	Sign       *secure.Sign
	Buf        *bytes.Buffer
	statusCode int
}

func (w *signWriter) Write(b []byte) (int, error) {
	w.Buf.Write(b)
	return w.Sign.Write(b)
}

func (w *signWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *signWriter) Close() error {
	hash := w.Sign.Sum()
	hashString := hex.EncodeToString(hash)

	w.ResponseWriter.Header().Set("HashSHA256", hashString)

	if w.statusCode > 0 {
		w.ResponseWriter.WriteHeader(w.statusCode)
	}

	w.ResponseWriter.Write(w.Buf.Bytes())

	if c, ok := w.ResponseWriter.(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("io.WriteCloser is unavailable on the writer")
}
