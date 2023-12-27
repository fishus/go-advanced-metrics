package handlers

import (
	"html/template"
	"net/http"

	"go.uber.org/zap"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	counters := storage.Counters()
	gauges := storage.Gauges()

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
		logger.Log.Error(err.Error(), zap.String("event", "list handler"), zap.Any("data", data))
	}
}
