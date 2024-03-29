package middleware

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"sync"

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
	once       sync.Once
}

func (w *signWriter) Write(b []byte) (int, error) {
	w.Sign.Write(b)
	return w.Buf.Write(b)
}

func (w *signWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *signWriter) Close() error {
	var err error

	w.once.Do(func() {
		hash := w.Sign.Sum()
		hashString := hex.EncodeToString(hash)
		w.ResponseWriter.Header().Set("HashSHA256", hashString)

		if w.statusCode > 0 {
			w.ResponseWriter.WriteHeader(w.statusCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		w.ResponseWriter.Write(w.Buf.Bytes())

		if c, ok := w.ResponseWriter.(io.WriteCloser); ok {
			err = c.Close()
			return
		}
		err = errors.New("io.WriteCloser is unavailable on the writer")
	})

	return err
}
