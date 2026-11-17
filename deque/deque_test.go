package deque_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/deque"
	"github.com/stretchr/testify/assert"
)

func TestConstructorsImplementDeque(t *testing.T) {
	assert.IsType(t, (*deque.ArrayDeque[int])(nil), deque.NewArray[int](0))
	assert.IsType(t, (*deque.ArrayDeque[int])(nil), deque.ArrayFrom([]int{1}))
	assert.IsType(t, (*deque.ArrayDeque[int])(nil), deque.ArrayClone([]int{1}))
	assert.IsType(t, (*deque.ArrayDeque[int])(nil), deque.ArrayFromSeq(slices.Values([]int{1})))
	assert.IsType(t, (*deque.LinkedDeque[int])(nil), deque.NewLinked[int]())
	assert.IsType(t, (*deque.LinkedDeque[int])(nil), deque.LinkedFrom([]int{1}))
	assert.IsType(t, (*deque.LinkedDeque[int])(nil), deque.LinkedFromSeq(slices.Values([]int{1})))
	var _ deque.Deque[int] = deque.NewArray[int](0)
	var _ deque.Deque[int] = deque.NewLinked[int]()
}

func TestConcreteZeroValues(t *testing.T) {
	var array deque.ArrayDeque[int]
	array.Append(1, 2)
	array.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, array.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &array))
	assert.Equal(t, []int{3, 4}, array.ToSlice())

	var linked deque.LinkedDeque[int]
	linked.Append(1, 2)
	linked.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, linked.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &linked))
	assert.Equal(t, []int{3, 4}, linked.ToSlice())
}

func TestArrayFrom(t *testing.T) {
	data := []int{1, 2}
	d := deque.ArrayFrom(data)

	assert.Equal(t, []int{1, 2}, d.ToSlice())
}

func TestArrayClone(t *testing.T) {
	data := []int{1, 2}
	d := deque.ArrayClone(data)

	_, _ = d.Shift()
	assert.Equal(t, []int{1, 2}, data)
	data[1] = 20
	assert.Equal(t, []int{2}, d.ToSlice())
}

func TestArrayFromSeq(t *testing.T) {
	d := deque.ArrayFromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, d.ToSlice())
}

func TestLinkedFrom(t *testing.T) {
	data := []int{1, 2}
	d := deque.LinkedFrom(data)
	data[0] = 10

	assert.Equal(t, []int{1, 2}, d.ToSlice())
}

func TestLinkedFromSeq(t *testing.T) {
	d := deque.LinkedFromSeq(slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, d.ToSlice())
}

func TestNewArray(t *testing.T) {
	assertDequeBehavior(t, deque.NewArray[int](1))
}

func TestNewLinked(t *testing.T) {
	assertDequeBehavior(t, deque.NewLinked[int]())
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() deque.Deque[int]{
		"ArrayFromNil":    func() deque.Deque[int] { return deque.ArrayFrom([]int(nil)) },
		"ArrayFromEmpty":  func() deque.Deque[int] { return deque.ArrayFrom(make([]int, 0, 4)) },
		"ArrayCloneNil":   func() deque.Deque[int] { return deque.ArrayClone([]int(nil)) },
		"ArrayCloneEmpty": func() deque.Deque[int] { return deque.ArrayClone(make([]int, 0, 4)) },
		"LinkedFromNil":   func() deque.Deque[int] { return deque.LinkedFrom([]int(nil)) },
		"LinkedFromEmpty": func() deque.Deque[int] { return deque.LinkedFrom(make([]int, 0, 4)) },
	} {
		t.Run(name, func(t *testing.T) {
			values := construct()
			assert.True(t, values.IsEmpty())
			values.Append(1)
			assert.Equal(t, []int{1}, values.ToSlice())
		})
	}
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, deque.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := deque.NewArray[int](0)
		mutable.Append(1, 2)

		view := deque.AsReadonly(mutable)
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, unmarshals)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Len())
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
