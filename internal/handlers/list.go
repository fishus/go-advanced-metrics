package handlers

import (
	"html/template"
	"net/http"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

// ListHandler processes the request GET /.
// Returns the values of all metrics.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	counters := config.Storage.CountersContext(r.Context())
	gauges := config.Storage.GaugesContext(r.Context())

	data := struct {
		Counters map[string]metrics.Counter
		Gauges   map[string]metrics.Gauge
	}{
		Counters: counters,
		Gauges:   gauges,
	}

	templates := template.Must(template.New("list.html").ParseFiles("templates/list.html"))
	err := templates.Execute(w, data)
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "list handler"), logger.Any("data", data))
	}
}
