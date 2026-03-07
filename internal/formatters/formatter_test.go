package formatters_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name       string
		collection *FakeCollection[int]
		format     string
		expected   string
	}{
		{
			name:       "Empty Value",
			collection: &FakeCollection[int]{data: []int{}},
			format:     "%v",
			expected:   "[]",
		},
		{
			name:       "Simple Value",
			collection: &FakeCollection[int]{data: []int{1, 2, 3}},
			format:     "%v",
			expected:   "[1 2 3]",
		},
		{
			name:       "Verbose",
			collection: &FakeCollection[int]{data: []int{1, 2}},
			format:     "%+v",
			expected:   "*formatters_test.FakeCollection[int]{len:2, cap:2} [1 2]",
		},
		{
			name:       "Go Syntax",
			collection: &FakeCollection[int]{data: []int{1, 2}},
			format:     "%#v",
			expected:   "*formatters_test.FakeCollection[int]{size:2, cap:2}",
		},
		{
			name: "Truncate",
			collection: &FakeCollection[int]{
				data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			},
			format:   "%v",
			expected: "[1 2 3 4 5 ...(+12 more)]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf(tt.format, tt.collection)
			assert.Equal(t, tt.expected, got)
		})
	}
}
