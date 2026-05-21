package deque_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/deque"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementDeque(t *testing.T) {
	var _ deque.Deque[int] = deque.NewArray[int](0)
	var _ deque.Deque[int] = deque.NewLinked[int]()
}

func TestNewArray(t *testing.T) {
	assertDequeBehavior(t, deque.NewArray[int](1))
}

func TestNewLinked(t *testing.T) {
	assertDequeBehavior(t, deque.NewLinked[int]())
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, deque.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := deque.NewArray[int](0)
		mutable.Append(1, 2)

		view := deque.AsReadonly[int](mutable)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, 2, view.Length())

		mutable.Prepend(0)
		assert.Equal(t, []int{0, 1, 2}, view.ToSlice())

		data, err := json.Marshal(view)
		assert.NoError(t, err)
		assert.JSONEq(t, `[0,1,2]`, string(data))
	})
}

func assertDequeBehavior(t *testing.T, d deque.Deque[int]) {
	t.Helper()

	assert.True(t, d.IsEmpty())

	d.Append(2, 3)
	d.Prepend(0, 1)

	assert.Equal(t, []int{0, 1, 2, 3}, d.ToSlice())
	assert.Equal(t, []int{0, 1, 2, 3}, slices.Collect(d.All()))

	front, ok := d.Front()
	assert.True(t, ok)
	assert.Equal(t, 0, front)

	back, ok := d.Back()
	assert.True(t, ok)
	assert.Equal(t, 3, back)

	x, ok := d.Shift()
	assert.True(t, ok)
	assert.Equal(t, 0, x)

	x, ok = d.Pop()
	assert.True(t, ok)
	assert.Equal(t, 3, x)

	data, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.JSONEq(t, `[1,2]`, string(data))

	err = json.Unmarshal([]byte(`[8,9]`), d)
	assert.NoError(t, err)
	assert.Equal(t, []int{8, 9}, d.ToSlice())

	d.Clear()
	assert.True(t, d.IsEmpty())
	assert.Nil(t, d.ToSlice())
}
