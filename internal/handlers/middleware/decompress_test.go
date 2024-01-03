package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestDecompress(t *testing.T) {
	r := chi.NewRouter()
	r.Use(Decompress)

	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testCases := []struct {
		name     string
		encoding string
		msg      []byte
		want     []byte
	}{
		{
			name:     "Positive: no content encoding",
			encoding: "",
			msg:      []byte(`{"message":"Hello"}`),
			want:     []byte(`{"message":"Hello"}`),
		},
		{
			name:     "Positive: gzip encoding",
			encoding: "gzip",
			msg:      gzipCompress(t, []byte(`{"message":"Hello"}`)),
			want:     []byte(`{"message":"Hello"}`),
		},
		{
			name:     "Positive: wrong content encoding",
			encoding: "none",
			msg:      []byte(`{"message":"Hello"}`),
			want:     []byte(`{"message":"Hello"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, tc.encoding, tc.msg)
			resp.Body.Close()
			assert.Equal(t, tc.want, respBody)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, encoding string, data []byte) (*http.Response, []byte) {
	req, err := http.NewRequest(http.MethodPost, (ts.URL + "/test"), bytes.NewBuffer(data))
	require.NoError(t, err)

	if encoding != "" {
		req.Header.Set("Content-Encoding", encoding)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	return resp, respBody
}

func gzipCompress(t *testing.T, data []byte) []byte {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	require.NoError(t, err)

	_, err = w.Write(data)
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)

	return b.Bytes()
}
