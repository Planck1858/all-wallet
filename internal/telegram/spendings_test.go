package telegram

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getExpenseArrFromStr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected []expanse
		withErr  bool
	}{
		{
			name:     "empty string",
			str:      "",
			expected: nil,
			withErr:  true,
		},
		{
			name: "success - 1 element, short time",
			str:  "+ 123.4 usd 10.06",
			expected: []expanse{
				{
					amount:   123.40,
					currency: "usd",
					date:     time.Date(time.Now().Year(), 6, 10, 0, 0, 0, 0, time.UTC),
				},
			},
			withErr: false,
		},
		{
			name: "success - 1 element, long time",
			str:  "+ 123.4 usd 10.06.2022",
			expected: []expanse{
				{
					amount:   123.40,
					currency: "usd",
					date:     time.Date(2022, 6, 10, 0, 0, 0, 0, time.UTC),
				},
			},
			withErr: false,
		},
		{
			name: "success - 1+ element, mixed time formats",
			str: `+ 123 usd 10.06.2022
- 123.40 usd 10.06
- 4321 eur 12.12.2012`,
			expected: []expanse{
				{
					amount:   123,
					currency: "usd",
					date:     time.Date(2022, 6, 10, 0, 0, 0, 0, time.UTC),
				},
				{
					amount:   -123.40,
					currency: "usd",
					date:     time.Date(time.Now().Year(), 6, 10, 0, 0, 0, 0, time.UTC),
				},
				{
					amount:   -4321,
					currency: "eur",
					date:     time.Date(2012, 12, 12, 0, 0, 0, 0, time.UTC),
				},
			},
			withErr: false,
		},
		{
			name: "success - 1 element, long time",
			str:  "++ 123.4 usd 10.06.2022",
			expected: []expanse{
				{
					amount:   123.40,
					currency: "usd",
					date:     time.Date(2022, 6, 10, 0, 0, 0, 0, time.UTC),
				},
			},
			withErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := getExpanseArrFromStr(tt.str)
			if err != nil {
				if tt.withErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
