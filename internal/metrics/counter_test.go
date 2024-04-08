package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCounter(t *testing.T) {
	type attr struct {
		name  string
		value int64
	}

	testCases := []struct {
		name    string
		attr    attr
		want    Counter
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			attr:    attr{name: "test", value: 10},
			want:    Counter{name: "test", value: 10},
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			attr:    attr{name: "", value: 10},
			want:    Counter{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			counter, err := NewCounter(tc.attr.name, tc.attr.value)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, tc.want, *counter)
		})
	}
}

func TestCounter_AddValue(t *testing.T) {
	testCases := []struct {
		name    string
		counter Counter
		value   int64
		want    int64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			value:   1,
			counter: Counter{"test", 2},
			want:    3,
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			value:   0,
			counter: Counter{"test", 1},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			value:   -1,
			counter: Counter{"test", 2},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.counter.AddValue(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, tc.counter.Value())
		})
	}
}
