package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type DecompressSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *DecompressSuite) SetupSuite() {
	r := chi.NewRouter()
	r.Use(Decompress)

	r.Post("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, r.Body)
		s.Require().NoError(err)
		defer r.Body.Close()
	})

	s.ts = httptest.NewServer(r)
	s.client = resty.New().SetBaseURL(s.ts.URL)
}

func (s *DecompressSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *DecompressSuite) sendRequest(encoding string, data []byte) *resty.Response {
	req := s.client.R().SetBody(data)

	if encoding != "" {
		req.SetHeader("Content-Encoding", encoding)
	}

	resp, err := req.Post("test/")
	s.Require().NoError(err)

	return resp
}

func (s *DecompressSuite) gzipCompress(data []byte) []byte {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	s.Require().NoError(err)

	_, err = w.Write(data)
	s.Require().NoError(err)

	err = w.Close()
	s.Require().NoError(err)

	return b.Bytes()
}

func (s *DecompressSuite) TestDecompress() {
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
			msg:      s.gzipCompress([]byte(`{"message":"Hello"}`)),
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
		s.Run(tc.name, func() {
			resp := s.sendRequest(tc.encoding, tc.msg)
			s.Equal(tc.want, resp.Body())
		})
	}
}

func TestDecompressSuite(t *testing.T) {
	suite.Run(t, new(DecompressSuite))
}
