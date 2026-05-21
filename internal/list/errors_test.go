package list

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIndexOutOfBoundError(t *testing.T) {
	err := NewIndexOutOfBoundError(-1, 3)
	message := "index -1 is out of bounds; maximum valid index is 3"
	template := "index %d is out of bounds; maximum valid index is %d"
	assert.EqualErrorf(t, err, message, template, -1, 3)
	assert.True(t, errors.Is(err, ErrIndexOutOfBound))
	assert.Equal(t, -1, err.Index())
	assert.Equal(t, 3, err.Limit())
}
