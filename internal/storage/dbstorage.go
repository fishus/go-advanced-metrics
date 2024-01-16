package storage

import (
	"context"
	"errors"
	"time"

	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type DBStorage struct {
	conn db.Connector
}

func NewDBStorage(conn db.Connector) *DBStorage {
	return &DBStorage{conn: conn}
}

func (ds *DBStorage) SetDBConn(conn db.Connector) {
	ds.conn = conn
}

// Gauge returns the gauge metric by name
func (ds *DBStorage) Gauge(name string) (metrics.Gauge, bool) {
	return ds.GaugeContext(context.Background(), name)
}

// GaugeContext returns the gauge metric by name
func (ds *DBStorage) GaugeContext(ctx context.Context, name string) (metrics.Gauge, bool) {
	if ds.conn == nil {
		return metrics.Gauge{}, false
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	row := ds.conn.QueryRow(ctxQuery, "SELECT value FROM metrics_gauge WHERE name = $1 LIMIT 1;", name)
	var value float64
	err := row.Scan(&value)
	if errors.Is(err, db.ErrNoRows) {
		return metrics.Gauge{}, false
	}
	if err != nil {
		logger.Log.Warn(err.Error())
		return metrics.Gauge{}, false
	}

	gauge, err := metrics.NewGauge(name, value)
	if err != nil {
		logger.Log.Warn(err.Error())
		return metrics.Gauge{}, false
	}

	return *gauge, true
}

// GaugeValue returns the gauge metric value by name
func (ds *DBStorage) GaugeValue(name string) (float64, bool) {
	return ds.GaugeValueContext(context.Background(), name)
}

// GaugeValueContext returns the gauge metric value by name
func (ds *DBStorage) GaugeValueContext(ctx context.Context, name string) (float64, bool) {
	if gauge, ok := ds.GaugeContext(ctx, name); ok {
		return gauge.Value(), ok
	}
	return 0, false
}

// Gauges returns all gauge metrics
func (ds *DBStorage) Gauges() map[string]metrics.Gauge {
	return ds.GaugesContext(context.Background())
}

// GaugesContext returns all gauge metrics
func (ds *DBStorage) GaugesContext(ctx context.Context) map[string]metrics.Gauge {
	gauges := map[string]metrics.Gauge{}

	if ds.conn == nil {
		return gauges
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	rows, err := ds.conn.Query(ctxQuery, "SELECT name, value FROM metrics_gauge;")
	if err != nil {
		logger.Log.Warn(err.Error())
		return gauges
	}
	defer rows.Close()

	for rows.Next() {
		var (
			gName  string
			gValue float64
		)

		err = rows.Scan(&gName, &gValue)
		if err != nil {
			logger.Log.Warn(err.Error())
			return map[string]metrics.Gauge{}
		}

		gauge, err := metrics.NewGauge(gName, gValue)
		if err != nil {
			logger.Log.Warn(err.Error())
			return map[string]metrics.Gauge{}
		}

		gauges[gName] = *gauge
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Warn(err.Error())
		return map[string]metrics.Gauge{}
	}

	return gauges
}

func (ds *DBStorage) SetGauge(name string, value float64) error {
	return ds.SetGaugeContext(context.Background(), name, value)
}

func (ds *DBStorage) SetGaugeContext(ctx context.Context, name string, value float64) error {
	if ds.conn == nil {
		return db.ErrNotConnected
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	res, err := ds.conn.Exec(ctxQuery, "INSERT INTO metrics_gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;", name, value)
	if err != nil {
		logger.Log.Warn(err.Error(), logger.String("status", res.String()))
		return err
	}

	return nil
}

// Counter returns the counter metric by name
func (ds *DBStorage) Counter(name string) (metrics.Counter, bool) {
	return ds.CounterContext(context.Background(), name)
}

// CounterContext returns the counter metric by name
func (ds *DBStorage) CounterContext(ctx context.Context, name string) (metrics.Counter, bool) {
	if ds.conn == nil {
		return metrics.Counter{}, false
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	row := ds.conn.QueryRow(ctxQuery, "SELECT value FROM metrics_counter WHERE name = $1 LIMIT 1;", name)
	var value int64
	err := row.Scan(&value)
	if errors.Is(err, db.ErrNoRows) {
		return metrics.Counter{}, false
	}
	if err != nil {
		logger.Log.Warn(err.Error())
		return metrics.Counter{}, false
	}

	counter, err := metrics.NewCounter(name, value)
	if err != nil {
		logger.Log.Warn(err.Error())
		return metrics.Counter{}, false
	}

	return *counter, true
}

// CounterValue returns the counter metric value by name
func (ds *DBStorage) CounterValue(name string) (int64, bool) {
	return ds.CounterValueContext(context.Background(), name)
}

// CounterValueContext returns the counter metric value by name
func (ds *DBStorage) CounterValueContext(ctx context.Context, name string) (int64, bool) {
	if counter, ok := ds.CounterContext(ctx, name); ok {
		return counter.Value(), ok
	}
	return 0, false
}

// Counters returns all counter metrics
func (ds *DBStorage) Counters() map[string]metrics.Counter {
	return ds.CountersContext(context.Background())
}

// CountersContext returns all counter metrics
func (ds *DBStorage) CountersContext(ctx context.Context) map[string]metrics.Counter {
	counters := map[string]metrics.Counter{}

	if ds.conn == nil {
		return counters
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	rows, err := ds.conn.Query(ctxQuery, "SELECT name, value FROM metrics_counter;")
	if err != nil {
		logger.Log.Warn(err.Error())
		return counters
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cName  string
			cValue int64
		)

		err = rows.Scan(&cName, &cValue)
		if err != nil {
			logger.Log.Warn(err.Error())
			return map[string]metrics.Counter{}
		}

		counter, err := metrics.NewCounter(cName, cValue)
		if err != nil {
			logger.Log.Warn(err.Error())
			return map[string]metrics.Counter{}
		}

		counters[cName] = *counter
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Warn(err.Error())
		return map[string]metrics.Counter{}
	}

	return counters
}

func (ds *DBStorage) AddCounter(name string, value int64) error {
	return ds.AddCounterContext(context.Background(), name, value)
}

func (ds *DBStorage) AddCounterContext(ctx context.Context, name string, value int64) error {
	if ds.conn == nil {
		return db.ErrNotConnected
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	res, err := ds.conn.Exec(ctxQuery, "INSERT INTO metrics_counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = metrics_counter.value + EXCLUDED.value;", name, value)
	if err != nil {
		logger.Log.Warn(err.Error(), logger.String("status", res.String()))
		return err
	}

	return nil
}

var _ MetricsStorager = (*DBStorage)(nil)
