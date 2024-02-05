package middleware

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/secure"
)

type SignSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *SignSuite) SetupSuite() {
	r := chi.NewRouter()
	r.Use(Sign([]byte("secret")))

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

func (s *SignSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *SignSuite) sendRequest(data []byte) *resty.Response {
	req := s.client.R().SetBody(data).
		SetHeader("Content-Type", "text/plain; charset=utf-8")

	resp, err := req.Post("test/")
	s.Require().NoError(err)

	return resp
}

func (s *SignSuite) TestSign() {
	testCases := []struct {
		name     string
		data     string
		wantHash func(data []byte) string
	}{
		{
			name: "Positive",
			data: `Аэрофотосъёмка ландшафта уже выявила земли богачей и процветающих крестьян.`,
			wantHash: func(data []byte) string {
				hash := secure.Hash(data, []byte("secret"))
				return hex.EncodeToString(hash[:])
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := []byte(tc.data)
			resp := s.sendRequest(data)
			respHashString := resp.Header().Get("HashSHA256")
			s.Equal(tc.wantHash(data), respHashString)
		})
	}
}

func TestSignSuite(t *testing.T) {
	suite.Run(t, new(SignSuite))
}
