package list

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIndexOutOfBoundsError(t *testing.T) {
	err := NewIndexOutOfBoundsError(-1, 3)
	assert.EqualError(t, err, "index -1 out of range for length 4")
	assert.True(t, errors.Is(err, ErrIndexOutOfBounds))
	var target *IndexOutOfBoundsError
	assert.True(t, errors.As(err, &target))
	assert.Same(t, err, target)
	assert.Equal(t, -1, err.Index())
	assert.Equal(t, 3, err.Limit())

	empty := NewIndexOutOfBoundsError(0, -1)
	assert.EqualError(t, empty, "index 0 out of range for length 0")
	assert.Equal(t, -1, empty.Limit())
}
