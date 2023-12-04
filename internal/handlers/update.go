package handlers

import (
	"fmt"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
	"net/http"
	"strconv"
	"strings"
)

var ms metrics.Storager = metrics.NewMemStorage()

// Обработка данных в формате /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf(`%s method not allowed`, r.Method), http.StatusMethodNotAllowed)
		return
	}

	uri, ok := strings.CutPrefix(r.RequestURI, `/update/`)
	if !ok {
		http.Error(w, `The path must start with /update/`, http.StatusBadRequest)
		return
	}

	uri = strings.TrimRight(uri, `/`)

	uriSlice := strings.Split(uri, `/`)

	var metricType, metricName string

	if len(uriSlice) == 0 || uriSlice[0] == "" {
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Metric type not specified`, http.StatusBadRequest)
		return
	}
	metricType = uriSlice[0]

	switch metricType {
	case "counter", "gauge":
		// При попытке передать запрос без имени метрики возвращать http.StatusNotFound.
		if len(uriSlice) < 2 || uriSlice[1] == "" {
			http.Error(w, `Metric name not specified`, http.StatusNotFound)
			return
		}
		metricName = uriSlice[1]

		if len(uriSlice) < 3 {
			http.Error(w, `Incorrect metric value`, http.StatusBadRequest)
			return
		}
	default:
		// При попытке передать запрос с некорректным типом метрики http.StatusBadRequest.
		http.Error(w, `Incorrect metric type`, http.StatusBadRequest)
		return
	}

	switch metricType {
	case "counter":
		var metricValue int64

		if i, err := strconv.ParseInt(uriSlice[2], 10, 64); err != nil || uriSlice[2] == "" {
			http.Error(w, `Incorrect metric value`, http.StatusBadRequest)
			return
		} else {
			metricValue = i
		}

		err := ms.AddCounter(metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "gauge":
		var metricValue float64

		if i, err := strconv.ParseFloat(uriSlice[2], 64); err != nil || uriSlice[2] == "" {
			http.Error(w, `Incorrect metric value`, http.StatusBadRequest)
			return
		} else {
			metricValue = i
		}

		err := ms.SetGauge(metricName, metricValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// При успешном приёме возвращать http.StatusOK.
	w.WriteHeader(http.StatusOK)
}
