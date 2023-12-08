package handlers

import "github.com/go-chi/chi/v5"

func ServerRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", UpdateHandler)
	return r
}
