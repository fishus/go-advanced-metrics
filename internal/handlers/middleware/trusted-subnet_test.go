package middleware

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type TrustedSubnetSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
	subnet *net.IPNet
}

func (s *TrustedSubnetSuite) SetupSuite() {
	_, subnet, err := net.ParseCIDR("192.168.0.0/24")
	s.Require().NoError(err)
	s.subnet = subnet

	r := chi.NewRouter()
	r.Use(TrustedSubnet(subnet))

	r.Post("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := io.Copy(w, r.Body)
		s.Require().NoError(err)
	})

	s.ts = httptest.NewServer(r)
	s.client = resty.New().SetBaseURL(s.ts.URL)
}

func (s *TrustedSubnetSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *TrustedSubnetSuite) sendRequest(ip string) *resty.Response {
	req := s.client.R().SetHeader("X-Real-IP", ip)

	resp, err := req.Post("test/")
	s.Require().NoError(err)

	return resp
}

func (s *TrustedSubnetSuite) TestTrustedSubnet() {
	testCases := []struct {
		name string
		ip   string
		code int
	}{
		{
			name: "Positive case",
			ip:   "192.168.0.123",
			code: 200,
		},
		{
			name: "Negative case",
			ip:   "127.0.0.1",
			code: 403,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp := s.sendRequest(tc.ip)
			s.Equal(tc.code, resp.StatusCode())
		})
	}
}

func TestTrustedSubnetSuite(t *testing.T) {
	suite.Run(t, new(TrustedSubnetSuite))
}
