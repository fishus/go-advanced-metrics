package handlers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ListHandlerSuite struct {
	suite.Suite
	ts     *httptest.Server
	client *resty.Client
}

func (s *ListHandlerSuite) SetupSuite() {
	s.ts = httptest.NewServer(ServerRouter())
	s.client = resty.New()

	// When we run tests, the current directory is always the folder containing the test file.
	// So we need to change the working directory to the app's root dir
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
}

func (s *ListHandlerSuite) TearDownSuite() {
	s.ts.Close()
}

func (s *ListHandlerSuite) SetupTest() {
	storage = metrics.NewMemStorage()
}

func (s *ListHandlerSuite) requestValue() *resty.Response {
	resp, err := s.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		Get(s.ts.URL + "/")
	s.Require().NoError(err)
	return resp
}

func (s *ListHandlerSuite) TestNoMetrics() {
	resp := s.requestValue()
	body := resp.RawBody()
	defer body.Close()

	s.Equal(http.StatusOK, resp.StatusCode())

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic(err)
	}

	s.Run("Counters", func() {
		s.Equal(0, doc.Find(`[data-id="counters"] li`).Length())
		s.Equal(1, doc.Find(`[data-id="no-counters"]`).Length())
	})

	s.Run("Gauges", func() {
		s.Equal(0, doc.Find(`[data-id="gauges"] li`).Length())
		s.Equal(1, doc.Find(`[data-id="no-gauges"]`).Length())
	})
}

func (s *ListHandlerSuite) TestCounterOnly() {
	_ = storage.AddCounter("a", 5)

	resp := s.requestValue()
	body := resp.RawBody()
	defer body.Close()

	s.Equal(http.StatusOK, resp.StatusCode())

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic(err)
	}

	s.Run("Counters", func() {
		s.Equal(1, doc.Find(`[data-id="counters"] li`).Length())

		doc.Find(`[data-id="counters"] li`).Each(func(i int, sel *goquery.Selection) {
			name := sel.Find("li strong").Text()
			value := sel.Find("li span").Text()
			s.Equal("a", name)
			s.Equal("5", value)
		})
	})

	s.Run("Gauges", func() {
		s.Equal(0, doc.Find(`[data-id="gauges"] li`).Length())
		s.Equal(1, doc.Find(`[data-id="no-gauges"]`).Length())
	})
}

func (s *ListHandlerSuite) TestGaugeOnly() {
	_ = storage.SetGauge("a", 12.57)

	resp := s.requestValue()
	body := resp.RawBody()
	defer body.Close()

	s.Equal(http.StatusOK, resp.StatusCode())

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic(err)
	}

	s.Run("Counters", func() {
		s.Equal(0, doc.Find(`[data-id="counters"] li`).Length())
		s.Equal(1, doc.Find(`[data-id="no-counters"]`).Length())
	})

	s.Run("Gauges", func() {
		s.Equal(1, doc.Find(`[data-id="gauges"] li`).Length())

		doc.Find(`[data-id="gauges"] li`).Each(func(i int, sel *goquery.Selection) {
			name := sel.Find("li strong").Text()
			value := sel.Find("li span").Text()
			s.Equal("a", name)
			s.Equal("12.57", value)
		})
	})
}

func (s *ListHandlerSuite) TestCounterAndGauge() {
	_ = storage.AddCounter("a", 5)
	_ = storage.SetGauge("b", 12.57)

	resp := s.requestValue()
	body := resp.RawBody()
	defer body.Close()

	s.Equal(http.StatusOK, resp.StatusCode())

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		panic(err)
	}

	s.Run("Counters", func() {
		s.Equal(1, doc.Find(`[data-id="counters"] li`).Length())

		doc.Find(`[data-id="counters"] li`).Each(func(i int, sel *goquery.Selection) {
			name := sel.Find("li strong").Text()
			value := sel.Find("li span").Text()
			s.Equal("a", name)
			s.Equal("5", value)
		})
	})

	s.Run("Gauges", func() {
		s.Equal(1, doc.Find(`[data-id="gauges"] li`).Length())

		doc.Find(`[data-id="gauges"] li`).Each(func(i int, sel *goquery.Selection) {
			name := sel.Find("li strong").Text()
			value := sel.Find("li span").Text()
			s.Equal("b", name)
			s.Equal("12.57", value)
		})
	})
}

func TestListHandlerSuite(t *testing.T) {
	suite.Run(t, new(ListHandlerSuite))
}
