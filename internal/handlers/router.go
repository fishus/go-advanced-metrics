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
	r.Use(middleware.RequestLogger(&logger.LogFormatter{}))

	r.Post("/update/", UpdateMetricsHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", UpdateMetricHandler)
	r.Post("/value/", ValueMetricsHandler)
	r.Get("/value/{metricType}/{metricName}", ValueMetricHandler)
	r.Get("/", ListHandler)
	return r
}
