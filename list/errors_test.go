package list_test

import (
	"errors"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stretchr/testify/assert"
)

func TestIndexOutOfBoundsError(t *testing.T) {
	l := list.NewArray[int](0)

	_, err := l.Get(0)

	assert.True(t, errors.Is(err, list.ErrIndexOutOfBounds))
	assert.EqualError(t, err, "cannot get index 0 out of range for length 0")

	var target *list.IndexOutOfBoundsError
	assert.True(t, errors.As(err, &target))
	assert.Equal(t, 0, target.Index())
	assert.Equal(t, -1, target.Limit())

	err = l.Insert(1, 7)
	assert.EqualError(t, err, "cannot insert at index 1 out of range for length 0")
	assert.ErrorIs(t, err, list.ErrIndexOutOfBounds)
	assert.True(t, errors.As(err, &target))
	assert.Equal(t, 0, target.Limit())

	l.Append(7)
	_, err = l.Get(1)
	assert.EqualError(t, err, "cannot get index 1 out of range for length 1")
	assert.ErrorIs(t, err, list.ErrIndexOutOfBounds)
}
