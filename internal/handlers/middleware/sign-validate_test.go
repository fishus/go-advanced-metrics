package middleware

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/secure"
)

type ValidateSignSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *ValidateSignSuite) SetupSuite() {
	r := chi.NewRouter()
	r.Use(ValidateSign([]byte("secret")))

	r.Post("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		defer r.Body.Close()
	})

	s.ts = httptest.NewServer(r)
	s.client = resty.New().SetBaseURL(s.ts.URL)
}

func (s *ValidateSignSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *ValidateSignSuite) sendRequest(hashString string, data []byte) *resty.Response {
	req := s.client.R().SetBody(data).
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		SetHeader("HashSHA256", hashString)

	resp, err := req.Post("test/")
	s.Require().NoError(err)

	return resp
}

func (s *ValidateSignSuite) TestValidateSign() {
	testCases := []struct {
		name       string
		hashString func(data []byte) string
		wantCode   int
	}{
		{
			name: "Positive",
			hashString: func(data []byte) string {
				hash := secure.Hash(data, []byte("secret"))
				return hex.EncodeToString(hash[:])
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Negative: wrong secret",
			hashString: func(data []byte) string {
				hash := secure.Hash(data, []byte("key"))
				return hex.EncodeToString(hash[:])
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Negative: wrong hash",
			hashString: func(data []byte) string {
				return "12345abc"
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Negative: empty hash",
			hashString: func(data []byte) string {
				return ""
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := []byte(`Аэрофотосъёмка ландшафта уже выявила земли богачей и процветающих крестьян.`)
			resp := s.sendRequest(tc.hashString(data), data)
			s.Equal(tc.wantCode, resp.StatusCode())
		})
	}
}

func TestValidateSignSuite(t *testing.T) {
	suite.Run(t, new(ValidateSignSuite))
}
