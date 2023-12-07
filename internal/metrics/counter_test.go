package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCounter_Add(t *testing.T) {
	testCases := []struct {
		name    string
		value   int64
		counter Counter
		want    int64
		wantErr bool
	}{
		{
			name:    "Positive case #1",
			value:   1,
			counter: Counter(2),
			want:    3,
			wantErr: false,
		},
		{
			name:    "Positive case #2",
			value:   0,
			counter: Counter(1),
			want:    1,
			wantErr: false,
		},
		{
			name:    "Negative case #1",
			value:   -1,
			counter: Counter(2),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.counter.Add(tc.value)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, int64(tc.counter))
		})
	}
}
