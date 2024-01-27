package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/fishus/go-advanced-metrics/internal/handlers/middleware"
	"github.com/fishus/go-advanced-metrics/internal/logger"
)

func ServerRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(mw.Decompress)
	r.Use(middleware.Compress(9, "application/json", "text/html"))
	r.Use(middleware.RequestLogger(&logger.LogFormatter{}))

	if secretKey != "" {
		r.Use(mw.ValidateSign([]byte(secretKey)))
		r.Use(mw.Sign([]byte(secretKey)))
	}

	r.Post("/update/", UpdateMetricsHandler)
	r.Post("/updates/", UpdatesMetricsHandler)
	r.Post("/update/{metricType}/{metricID}/{metricValue}", UpdateMetricHandler)
	r.Post("/value/", ValueMetricsHandler)
	r.Get("/value/{metricType}/{metricID}", ValueMetricHandler)
	r.Get("/", ListHandler)
	r.Get("/ping", PingDBHandler)
	return r
}
