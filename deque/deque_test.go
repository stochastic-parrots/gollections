package deque_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/deque"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementDeque(t *testing.T) {
	var _ *deque.ArrayDeque[int] = deque.Array[int]().New(0)
	var _ *deque.ArrayDeque[int] = deque.Array[int]().From([]int{1})
	var _ *deque.ArrayDeque[int] = deque.Array[int]().Clone([]int{1})
	var _ *deque.ArrayDeque[int] = deque.Array[int]().FromSeq(slices.Values([]int{1}))
	var _ *deque.LinkedDeque[int] = deque.Linked[int]().New()
	var _ *deque.LinkedDeque[int] = deque.Linked[int]().From([]int{1})
	var _ *deque.LinkedDeque[int] = deque.Linked[int]().FromSeq(slices.Values([]int{1}))
	var _ deque.Deque[int] = deque.Array[int]().New(0)
	var _ deque.Deque[int] = deque.Linked[int]().New()
}

func TestConcreteZeroValues(t *testing.T) {
	var array deque.ArrayDeque[int]
	array.Append(1, 2)
	array.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, array.ToSlice())

	var linked deque.LinkedDeque[int]
	linked.Append(1, 2)
	linked.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, linked.ToSlice())
}

func TestArrayFactory_From(t *testing.T) {
	data := []int{1, 2}
	deque := deque.Array[int]().From(data)

	_, _ = deque.Shift()
	assert.Equal(t, []int{0, 2}, data)
}

func TestArrayFactory_Clone(t *testing.T) {
	data := []int{1, 2}
	deque := deque.Array[int]().Clone(data)

	_, _ = deque.Shift()
	assert.Equal(t, []int{1, 2}, data)
}

func TestArrayFactory_FromSeq(t *testing.T) {
	deque := deque.Array[int]().FromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestLinkedFactory_From(t *testing.T) {
	data := []int{1, 2}
	deque := deque.Linked[int]().From(data)
	data[0] = 10

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestLinkedFactory_FromSeq(t *testing.T) {
	deque := deque.Linked[int]().FromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestArrayFactory_New(t *testing.T) {
	assertDequeBehavior(t, deque.Array[int]().New(1))
}

func TestLinkedFactory_New(t *testing.T) {
	assertDequeBehavior(t, deque.Linked[int]().New())
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, deque.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := deque.Array[int]().New(0)
		mutable.Append(1, 2)

		view := deque.AsReadonly[int](mutable)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Length())
		assert.False(t, view.IsEmpty())

		front, ok := view.Front()
		assert.True(t, ok)
		assert.Equal(t, 1, front)

		back, ok := view.Back()
		assert.True(t, ok)
		assert.Equal(t, 2, back)

		assert.Equal(t, "[1 2]", view.String())

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

func collectIndexes(seq func(func(int, int) bool)) []int {
	var indexes []int
	for idx := range seq {
		indexes = append(indexes, idx)
	}
	return indexes
}
