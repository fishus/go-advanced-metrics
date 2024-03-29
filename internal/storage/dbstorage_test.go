package storage

import (
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/suite"

	"github.com/fishus/go-advanced-metrics/internal/logger"
	"github.com/fishus/go-advanced-metrics/internal/metrics"
)

type DBStorageSuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
}

func (s *DBStorageSuite) SetupSuite() {
	err := logger.Initialize("debug")
	s.Require().NoError(err)

	mock, err := pgxmock.NewPool()
	s.Require().NoError(err)
	s.mock = mock
}

func (s *DBStorageSuite) TearDownSuite() {
	logger.Log.Sync()
	s.mock.Close()
}

func (s *DBStorageSuite) TearDownTest() {
	s.mock.Reset()
}

func (s *DBStorageSuite) TestNewDBStorage() {
	want := &DBStorage{}
	got := NewDBStorage(nil)
	s.Equal(want, got)
}

func (s *DBStorageSuite) TestGauge() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("a").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(float64(2.1)))
	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("b").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(float64(-1.5)))
	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("c").WillReturnRows(s.mock.NewRows([]string{"value"}))

	type want struct {
		gauge metrics.Gauge
		ok    bool
	}

	testCases := []struct {
		name string
		key  string
		want want
	}{
		{
			name: "Positive case #1",
			key:  "a",
			want: want{
				gauge: func() metrics.Gauge {
					g, _ := metrics.NewGauge("a", 2.1)
					return *g
				}(),
				ok: true,
			},
		},
		{
			name: "Positive case #2",
			key:  "b",
			want: want{
				gauge: func() metrics.Gauge {
					g, _ := metrics.NewGauge("b", -1.5)
					return *g
				}(),
				ok: true,
			},
		},
		{
			name: "Negative case #1",
			key:  "c",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			gauge, ok := ds.Gauge(tc.key)

			s.Require().Equal(tc.want.ok, ok)
			if tc.want.ok {
				s.EqualValues(tc.want.gauge, gauge)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestGaugeValue() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("a").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(float64(2.1)))
	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("b").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(float64(-1.5)))
	s.mock.ExpectQuery("^SELECT value FROM metrics_gauge WHERE (.+) LIMIT 1;$").WithArgs("c").WillReturnRows(s.mock.NewRows([]string{"value"}))

	type want struct {
		value float64
		ok    bool
	}

	testCases := []struct {
		name string
		key  string
		want want
	}{
		{
			name: "Positive case #1",
			key:  "a",
			want: want{
				value: 2.1,
				ok:    true,
			},
		},
		{
			name: "Positive case #2",
			key:  "b",
			want: want{
				value: -1.5,
				ok:    true,
			},
		},
		{
			name: "Negative case #1",
			key:  "c",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			g, ok := ds.GaugeValue(tc.key)

			s.Require().Equal(tc.want.ok, ok)
			if tc.want.ok {
				s.EqualValues(tc.want.value, g)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestGauges() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("SELECT name, value FROM metrics_gauge;").
		WillReturnRows(s.mock.NewRows([]string{"name", "value"}).
			AddRow("b", float64(2.1)).
			AddRow("a", float64(1.0)))

	want := map[string]metrics.Gauge{}
	a, _ := metrics.NewGauge("a", 1.0)
	b, _ := metrics.NewGauge("b", 2.1)
	want["a"] = *a
	want["b"] = *b

	s.Equal(want, ds.Gauges())

	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestGaugesFiltered() {
	ds := NewDBStorage(s.mock)

	filter := []string{"a", "b"}

	s.mock.ExpectQuery("SELECT name, value FROM metrics_gauge WHERE (.+);").
		WithArgs(filter).
		WillReturnRows(s.mock.NewRows([]string{"name", "value"}).
			AddRow("b", float64(2.1)).
			AddRow("a", float64(1.0)))

	want := map[string]metrics.Gauge{}
	a, _ := metrics.NewGauge("a", 1.0)
	b, _ := metrics.NewGauge("b", 2.1)
	want["a"] = *a
	want["b"] = *b

	s.Equal(want, ds.Gauges(FilterNames(filter)))

	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestSetGauge() {
	ds := NewDBStorage(s.mock)

	insertSQL := `^INSERT INTO metrics_gauge (.+) VALUES (.+) ON CONFLICT \(name\) DO UPDATE SET value = EXCLUDED.value\;$`

	s.mock.ExpectExec(insertSQL).WithArgs("a", float64(5.0)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	s.mock.ExpectExec(insertSQL).WithArgs("a", float64(-5.0)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	s.mock.ExpectExec(insertSQL).WithArgs("b", float64(3)).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	testCases := []struct {
		name    string
		key     string
		value   float64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			key:     "a",
			value:   5.0,
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			key:     "a",
			value:   -5.0,
			wantErr: false,
		},
		{
			name:    "Positive case #3",
			key:     "b",
			value:   3,
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			key:     "",
			value:   5.0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := ds.SetGauge(tc.key, tc.value)
			if tc.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestResetGauges() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectExec(`^TRUNCATE metrics_gauge\;$`).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := ds.ResetGauges()
	s.Require().NoError(err)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestCounter() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("^SELECT (.+) FROM metrics_counter WHERE (.+)$").WithArgs("a").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(int64(21)))
	s.mock.ExpectQuery("^SELECT (.+) FROM metrics_counter WHERE (.+)$").WithArgs("b").WillReturnRows(s.mock.NewRows([]string{"value"}))

	type want struct {
		counter metrics.Counter
		ok      bool
	}

	testCases := []struct {
		name string
		key  string
		want want
	}{
		{
			name: "Positive case #1",
			key:  "a",
			want: want{
				counter: func() metrics.Counter {
					c, _ := metrics.NewCounter("a", 21)
					return *c
				}(),
				ok: true,
			},
		},
		{
			name: "Negative case #1",
			key:  "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			counter, ok := ds.Counter(tc.key)

			s.Require().Equal(tc.want.ok, ok)
			if tc.want.ok {
				s.EqualValues(tc.want.counter, counter)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestCounterValue() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("^SELECT (.+) FROM metrics_counter WHERE (.+)$").WithArgs("a").WillReturnRows(s.mock.NewRows([]string{"value"}).AddRow(int64(21)))
	s.mock.ExpectQuery("^SELECT (.+) FROM metrics_counter WHERE (.+)$").WithArgs("b").WillReturnRows(s.mock.NewRows([]string{"value"}))

	type want struct {
		value int64
		ok    bool
	}

	testCases := []struct {
		name string
		key  string
		want want
	}{
		{
			name: "Positive case #1",
			key:  "a",
			want: want{
				value: 21,
				ok:    true,
			},
		},
		{
			name: "Negative case #1",
			key:  "b",
			want: want{
				ok: false,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			c, ok := ds.CounterValue(tc.key)

			s.Require().Equal(tc.want.ok, ok)
			if tc.want.ok {
				s.EqualValues(tc.want.value, c)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestCounters() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectQuery("SELECT name, value FROM metrics_counter;").
		WillReturnRows(s.mock.NewRows([]string{"name", "value"}).
			AddRow("a", int64(1)).
			AddRow("b", int64(100)))

	want := map[string]metrics.Counter{}
	a, _ := metrics.NewCounter("a", 1)
	b, _ := metrics.NewCounter("b", 100)
	want["a"] = *a
	want["b"] = *b

	s.Equal(want, ds.Counters())

	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestCountersFiltered() {
	ds := NewDBStorage(s.mock)

	filter := []string{"a", "b"}

	s.mock.ExpectQuery("SELECT name, value FROM metrics_counter WHERE (.+);").
		WithArgs(filter).
		WillReturnRows(s.mock.NewRows([]string{"name", "value"}).
			AddRow("a", int64(1)).
			AddRow("b", int64(100)))

	want := map[string]metrics.Counter{}
	a, _ := metrics.NewCounter("a", 1)
	b, _ := metrics.NewCounter("b", 100)
	want["a"] = *a
	want["b"] = *b

	s.Equal(want, ds.Counters(FilterNames(filter)))

	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestAddCounter() {
	ds := NewDBStorage(s.mock)

	insertSQL := `^INSERT INTO metrics_counter (.+) VALUES (.+) ON CONFLICT \(name\) DO UPDATE SET value \= metrics_counter\.value \+ EXCLUDED\.value\;$`

	s.mock.ExpectExec(insertSQL).WithArgs("a", int64(1)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	s.mock.ExpectExec(insertSQL).WithArgs("b", int64(2)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	testCases := []struct {
		name    string
		key     string
		value   int64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			key:     "a",
			value:   1,
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			key:     "b",
			value:   2,
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			key:     "a",
			value:   -1,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			err := ds.AddCounter(tc.key, tc.value)
			if tc.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestResetCounters() {
	ds := NewDBStorage(s.mock)

	s.mock.ExpectExec(`^TRUNCATE metrics_counter\;$`).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := ds.ResetCounters()
	s.Require().NoError(err)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBStorageSuite) TestInsertBatch() {
	type counter struct {
		name  string
		value int64
	}

	type gauge struct {
		name  string
		value float64
	}

	testCases := []struct {
		name     string
		counters []counter
		gauges   []gauge
	}{
		{
			name: "Insert counters",
			counters: []counter{
				{name: "a", value: 1},
				{name: "b", value: 2},
			},
		},
		{
			name: "Insert gauges",
			gauges: []gauge{
				{name: "a", value: 1.2},
				{name: "b", value: 2.3},
			},
		},
		{
			name: "Insert counters and gauges",
			counters: []counter{
				{name: "a", value: 3},
				{name: "b", value: 4},
			},
			gauges: []gauge{
				{name: "a", value: 3.4},
				{name: "b", value: 4.5},
			},
		},
		{
			name: "Insert nothing",
		},
	}

	ds := NewDBStorage(s.mock)
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if len(tc.counters) == 0 && len(tc.gauges) == 0 {
				err := ds.InsertBatch()
				s.Require().NoError(err)
				return
			}

			s.mock.ExpectBegin()

			var countersBatch []metrics.Counter
			if len(tc.counters) > 0 {
				s.mock.ExpectPrepare("insert-counter", `^INSERT INTO metrics_counter \(name, value\) VALUES (.+) ON CONFLICT \(name\) DO UPDATE SET value \= metrics_counter\.value \+ EXCLUDED\.value;$`)
				for _, v := range tc.counters {
					s.mock.ExpectExec("insert-counter").WithArgs(v.name, int64(v.value)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
					c, err := metrics.NewCounter(v.name, v.value)
					if err == nil {
						countersBatch = append(countersBatch, *c)
					}
				}
			}

			var gaugesBatch []metrics.Gauge
			if len(tc.gauges) > 0 {
				s.mock.ExpectPrepare("insert-gauge", `^INSERT INTO metrics_gauge \(name, value\) VALUES (.+) ON CONFLICT \(name\) DO UPDATE SET value \= EXCLUDED\.value;$`)
				for _, v := range tc.gauges {
					s.mock.ExpectExec("insert-gauge").WithArgs(v.name, float64(v.value)).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
					g, err := metrics.NewGauge(v.name, v.value)
					if err == nil {
						gaugesBatch = append(gaugesBatch, *g)
					}
				}
			}

			s.mock.ExpectCommit()

			var err error
			if len(tc.counters) > 0 && len(tc.gauges) > 0 {
				err = ds.InsertBatch(WithCounters(countersBatch), WithGauges(gaugesBatch))
			} else if len(tc.counters) > 0 {
				err = ds.InsertBatch(WithCounters(countersBatch))
			} else if len(tc.gauges) > 0 {
				err = ds.InsertBatch(WithGauges(gaugesBatch))
			} else {
				err = ds.InsertBatch()
			}

			s.Require().NoError(err)
		})
	}
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func TestDBStorageSuite(t *testing.T) {
	suite.Run(t, new(DBStorageSuite))
}
