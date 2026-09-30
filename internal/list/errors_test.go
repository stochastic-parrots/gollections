package list

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIndexOutOfBoundsError(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		index     int
		limit     int
		message   string
	}{
		{"Get", OperationGet, 0, -1, "cannot get index 0 out of range for length 0"},
		{"Set", OperationSet, -1, 3, "cannot set index -1 out of range for length 4"},
		{"Insert", OperationInsert, 1, 0, "cannot insert at index 1 out of range for length 0"},
		{"Remove", OperationRemove, 4, 3, "cannot remove index 4 out of range for length 4"},
		{"Replace", OperationReplace, 4, 3, "cannot replace index 4 out of range for length 4"},
		{"Unspecified", 0, -1, 3, "index -1 out of range for length 4"},
		{"Unknown", Operation(255), -1, 3, "index -1 out of range for length 4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewIndexOutOfBoundsError(tt.operation, tt.index, tt.limit)
			assert.EqualError(t, err, tt.message)
			assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
			var target *IndexOutOfBoundsError
			assert.True(t, errors.As(err, &target))
			assert.Same(t, err, target)
			assert.Equal(t, tt.index, err.Index())
			assert.Equal(t, tt.limit, err.Limit())
		})
	}
}

func TestErrIndexOutOfBounds(t *testing.T) {
	assert.EqualError(t, ErrIndexOutOfBounds, "index 0 out of range for length 0")
	assert.Equal(t, -1, ErrIndexOutOfBounds.Limit())
}
