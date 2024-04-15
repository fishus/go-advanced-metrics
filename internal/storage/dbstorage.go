package storage

import (
	"context"
	"errors"
	"os"
	"time"

	db "github.com/fishus/go-advanced-metrics/internal/database"
	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type DBStorage struct {
	pool db.Connector
}

func NewDBStorage(pool db.Connector) *DBStorage {
	return &DBStorage{pool: pool}
}

func (ds *DBStorage) SetDBPool(pool db.Connector) {
	ds.pool = pool
}

func (ds *DBStorage) GetDBPool() (db.Connector, error) {
	if ds.pool == nil {
		return nil, db.ErrNotConnected
	}

	// Delay after unsuccessful request
	retryDelay := []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
		0,
	}

	var err error
	for _, delay := range retryDelay {
		ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Second))
		defer cancel()
		err = ds.pool.Ping(ctx)

		if err == nil || !db.IsConnectionException(err) {
			return ds.pool, nil
		}
		time.Sleep(delay)
	}

	return nil, err
}

// Gauge returns the gauge metric by name
func (ds *DBStorage) Gauge(name string) (metrics.Gauge, bool) {
	return ds.GaugeContext(context.Background(), name)
}

// GaugeContext returns the gauge metric by name
func (ds *DBStorage) GaugeContext(ctx context.Context, name string) (metrics.Gauge, bool) {
	pool, err := ds.GetDBPool()
	if err != nil {
		return metrics.Gauge{}, false
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	row := pool.QueryRow(ctxQuery, "SELECT value FROM metrics_gauge WHERE name = $1 LIMIT 1;", name)
	var value float64
	err = row.Scan(&value)
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
func (ds *DBStorage) Gauges(filters ...StorageFilter) map[string]metrics.Gauge {
	return ds.GaugesContext(context.Background(), filters...)
}

// GaugesContext returns all gauge metrics
func (ds *DBStorage) GaugesContext(ctx context.Context, filters ...StorageFilter) map[string]metrics.Gauge {
	gauges := map[string]metrics.Gauge{}

	var (
		rows db.Rows
		err  error
	)
	pool, err := ds.GetDBPool()
	if err != nil {
		return gauges
	}

	f := &StorageFilters{}
	for _, filter := range filters {
		filter(f)
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	if len(f.names) > 0 {
		rows, err = pool.Query(ctxQuery, "SELECT name, value FROM metrics_gauge WHERE name = ANY($1);", f.names)
	} else {
		rows, err = pool.Query(ctxQuery, "SELECT name, value FROM metrics_gauge;")
	}
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

		if err = rows.Scan(&gName, &gValue); err != nil {
			logger.Log.Warn(err.Error())
			return map[string]metrics.Gauge{}
		}

		gauge, err2 := metrics.NewGauge(gName, gValue)
		if err2 != nil {
			logger.Log.Warn(err2.Error())
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
	pool, err := ds.GetDBPool()
	if err != nil {
		return err
	}

	if _, err = metrics.NewGauge(name, value); err != nil {
		return err
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	_, err = pool.Exec(ctxQuery, "INSERT INTO metrics_gauge (name, value) VALUES (@name, @value) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;",
		db.NamedArgs{"name": name, "value": value})
	if err != nil {
		return err
	}

	return nil
}

func (ds *DBStorage) ResetGauges() error {
	pool, err := ds.GetDBPool()
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = pool.Exec(ctx, "TRUNCATE metrics_gauge;")
	if err != nil {
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
	pool, err := ds.GetDBPool()
	if err != nil {
		return metrics.Counter{}, false
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	row := pool.QueryRow(ctxQuery, "SELECT value FROM metrics_counter WHERE name = $1 LIMIT 1;", name)
	var value int64
	err = row.Scan(&value)
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
func (ds *DBStorage) Counters(filters ...StorageFilter) map[string]metrics.Counter {
	return ds.CountersContext(context.Background(), filters...)
}

// CountersContext returns all counter metrics
func (ds *DBStorage) CountersContext(ctx context.Context, filters ...StorageFilter) map[string]metrics.Counter {
	counters := map[string]metrics.Counter{}

	var (
		rows db.Rows
		err  error
	)
	pool, err := ds.GetDBPool()
	if err != nil {
		return counters
	}

	f := &StorageFilters{}
	for _, filter := range filters {
		filter(f)
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	if len(f.names) > 0 {
		rows, err = pool.Query(ctxQuery, "SELECT name, value FROM metrics_counter WHERE name = ANY($1);", f.names)
	} else {
		rows, err = pool.Query(ctxQuery, "SELECT name, value FROM metrics_counter;")
	}
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

		counter, err2 := metrics.NewCounter(cName, cValue)
		if err2 != nil {
			logger.Log.Warn(err2.Error())
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
	pool, err := ds.GetDBPool()
	if err != nil {
		return err
	}

	if _, err = metrics.NewCounter(name, value); err != nil {
		return err
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (3 * time.Second))
	defer cancel()

	_, err = pool.Exec(ctxQuery, "INSERT INTO metrics_counter (name, value) VALUES (@name, @value) ON CONFLICT (name) DO UPDATE SET value = metrics_counter.value + EXCLUDED.value;",
		db.NamedArgs{"name": name, "value": value})
	if err != nil {
		return err
	}

	return nil
}

func (ds *DBStorage) ResetCounters() error {
	pool, err := ds.GetDBPool()
	if err != nil {
		return err
	}

	ctx := context.Background()

	_, err = pool.Exec(ctx, "TRUNCATE metrics_counter;")
	if err != nil {
		return err
	}

	return nil
}

// MigrateCreateSchema Создать все необходимые таблицы в базе данных.
func (ds *DBStorage) MigrateCreateSchema(ctx context.Context) {
	pool, err := ds.GetDBPool()
	if err != nil {
		logger.Log.Warn(db.ErrNotConnected.Error())
		return
	}

	dump, err := os.ReadFile(`db/migration/create_schema.sql`)
	if err != nil {
		logger.Log.Warn(err.Error())
		return
	}

	ctxQuery, cancel := context.WithTimeout(ctx, (30 * time.Second))
	defer cancel()

	res, err := pool.Exec(ctxQuery, string(dump))
	if err != nil {
		logger.Log.Warn(err.Error(), logger.String("status", res.String()))
		return
	}
}

func (ds *DBStorage) Reset() error {
	gErr := ds.ResetGauges()
	cErr := ds.ResetCounters()
	return errors.Join(gErr, cErr)
}

func (ds *DBStorage) InsertBatch(opts ...StorageOption) error {
	return ds.InsertBatchContext(context.Background(), opts...)
}

func (ds *DBStorage) InsertBatchContext(ctx context.Context, opts ...StorageOption) error {
	pool, err := ds.GetDBPool()
	if err != nil {
		return err
	}

	o := &StorageOptions{}
	for _, opt := range opts {
		opt(o)
	}

	if len(o.gauges) == 0 && len(o.counters) == 0 {
		return nil
	}

	ctxTx, cancel := context.WithTimeout(ctx, (30 * time.Second))
	defer cancel()

	tx, err := pool.Begin(ctxTx)
	if err != nil {
		return err
	}

	if len(o.counters) > 0 {
		ctxPrepareCounter, cancel := context.WithTimeout(ctxTx, (3 * time.Second))
		defer cancel()

		stmtCounter, err := tx.Prepare(ctxPrepareCounter, "insert-counter",
			"INSERT INTO metrics_counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = metrics_counter.value + EXCLUDED.value;")
		if err != nil {
			if errR := tx.Rollback(ctxTx); errR != nil {
				return errors.Join(err, errR)
			}
			return err
		}

		for _, counter := range o.counters {
			ctxQuery, cancel := context.WithTimeout(ctxTx, (3 * time.Second))
			defer cancel()

			_, err := tx.Exec(ctxQuery, stmtCounter.Name, counter.Name(), counter.Value())
			if err != nil {
				if errR := tx.Rollback(ctxTx); errR != nil {
					return errors.Join(err, errR)
				}
				return err
			}
		}
	}

	if len(o.gauges) > 0 {
		ctxPrepareGauge, cancel := context.WithTimeout(ctxTx, (3 * time.Second))
		defer cancel()

		stmtGauge, err := tx.Prepare(ctxPrepareGauge, "insert-gauge",
			"INSERT INTO metrics_gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;")
		if err != nil {
			if errR := tx.Rollback(ctxTx); errR != nil {
				return errors.Join(err, errR)
			}
			return err
		}

		for _, gauge := range o.gauges {
			ctxQuery, cancel := context.WithTimeout(ctxTx, (3 * time.Second))
			defer cancel()

			_, err := tx.Exec(ctxQuery, stmtGauge.Name, gauge.Name(), gauge.Value())
			if err != nil {
				if errR := tx.Rollback(ctxTx); errR != nil {
					return errors.Join(err, errR)
				}
				return err
			}
		}
	}

	tx.Commit(ctxTx)

	return nil
}

var _ MetricsStorager = (*DBStorage)(nil)
