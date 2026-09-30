package list_test

import (
	"errors"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stretchr/testify/assert"
)

func TestIndexOutOfBoundError(t *testing.T) {
	l := list.NewArray[int](0)

	_, err := l.Get(0)

	assert.True(t, errors.Is(err, list.ErrIndexOutOfBound))

	var target *list.IndexOutOfBoundError
	assert.True(t, errors.As(err, &target))
	assert.Equal(t, 0, target.Index())
	assert.Equal(t, -1, target.Limit())
}
