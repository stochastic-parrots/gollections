package deque_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/deque"
	"github.com/stretchr/testify/assert"
)

var _ json.Marshaler = (*deque.ArrayDeque[int])(nil)
var _ fmt.Stringer = (*deque.ArrayDeque[int])(nil)
var _ fmt.Formatter = (*deque.ArrayDeque[int])(nil)
var _ json.Unmarshaler = (*deque.ArrayDeque[int])(nil)

var _ json.Marshaler = (*deque.LinkedDeque[int])(nil)
var _ fmt.Stringer = (*deque.LinkedDeque[int])(nil)
var _ fmt.Formatter = (*deque.LinkedDeque[int])(nil)
var _ json.Unmarshaler = (*deque.LinkedDeque[int])(nil)

func TestConstructorsImplementDeque(t *testing.T) {
	var _ *deque.ArrayDeque[int] = deque.NewArray[int](0)
	var _ *deque.ArrayDeque[int] = deque.ArrayFrom[int]([]int{1})
	var _ *deque.ArrayDeque[int] = deque.ArrayClone[int]([]int{1})
	var _ *deque.ArrayDeque[int] = deque.ArrayFromSeq[int](slices.Values([]int{1}))
	var _ *deque.LinkedDeque[int] = deque.NewLinked[int]()
	var _ *deque.LinkedDeque[int] = deque.LinkedFrom[int]([]int{1})
	var _ *deque.LinkedDeque[int] = deque.LinkedFromSeq[int](slices.Values([]int{1}))
	var _ deque.Deque[int] = deque.NewArray[int](0)
	var _ deque.Deque[int] = deque.NewLinked[int]()
}

func TestConcreteZeroValues(t *testing.T) {
	var array deque.ArrayDeque[int]
	array.Appends(1, 2)
	array.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, array.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &array))
	assert.Equal(t, []int{3, 4}, array.ToSlice())

	var linked deque.LinkedDeque[int]
	linked.Appends(1, 2)
	linked.Prepend(0)
	assert.Equal(t, []int{0, 1, 2}, linked.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &linked))
	assert.Equal(t, []int{3, 4}, linked.ToSlice())
}

func TestArrayFrom(t *testing.T) {
	data := []int{1, 2}
	deque := deque.ArrayFrom[int](data)

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestArrayClone(t *testing.T) {
	data := []int{1, 2}
	deque := deque.ArrayClone[int](data)

	_, _ = deque.Shift()
	assert.Equal(t, []int{1, 2}, data)
	data[1] = 20
	assert.Equal(t, []int{2}, deque.ToSlice())
}

func TestArrayFromSeq(t *testing.T) {
	deque := deque.ArrayFromSeq[int](slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestLinkedFrom(t *testing.T) {
	data := []int{1, 2}
	deque := deque.LinkedFrom[int](data)
	data[0] = 10

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
}

func TestLinkedFromSeq(t *testing.T) {
	deque := deque.LinkedFromSeq[int](slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, deque.ToSlice())
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
		mutable.Appends(1, 2)

		view := deque.AsReadonly[int](mutable)
		assert.NotImplements(t, (*deque.Deque[int])(nil), view)
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

func TestReadonly_String(t *testing.T) {
	for source, construct := range readonlySources() {
		t.Run(source, func(t *testing.T) {
			for _, test := range []struct {
				name   string
				values []int
				want   string
			}{
				{"Empty", nil, "[]"},
				{"Values", []int{1, 2, 3}, "[1 2 3]"},
				{"DisplayLimit", []int{1, 2, 3, 4, 5}, "[1 2 3 4 5]"},
				{"Truncated", []int{1, 2, 3, 4, 5, 6}, "[1 2 3 4 5 ...(+1 more)]"},
			} {
				t.Run(test.name, func(t *testing.T) {
					values := construct(slices.Clone(test.values))
					view := deque.AsReadonly[int](values)
					assert.Equal(t, test.want, view.String())
					assert.Equal(t, test.want, fmt.Sprint(view))
				})
			}
		})
	}
}

func TestReadonly_MarshalJSON(t *testing.T) {
	for source, construct := range readonlySources() {
		t.Run(source, func(t *testing.T) {
			for _, test := range []struct {
				name   string
				values []int
			}{
				{"Empty", nil},
				{"Values", []int{1, 2, 3}},
				{"BeyondDisplayLimit", []int{1, 2, 3, 4, 5, 6}},
			} {
				t.Run(test.name, func(t *testing.T) {
					values := construct(slices.Clone(test.values))
					view := deque.AsReadonly[int](values)
					data, err := view.MarshalJSON()
					assert.NoError(t, err)
					if len(test.values) == 0 {
						assert.Equal(t, "[]", string(data))
					} else {
						var decoded []int
						assert.NoError(t, json.Unmarshal(data, &decoded))
						assert.Equal(t, test.values, decoded)
					}
				})
			}
		})
	}

	t.Run("ElementError", func(t *testing.T) {
		for name, values := range map[string]deque.Deque[chan int]{
			"Array":  deque.ArrayFrom([]chan int{make(chan int)}),
			"Linked": deque.LinkedFrom([]chan int{make(chan int)}),
		} {
			t.Run(name, func(t *testing.T) {
				view := deque.AsReadonly[chan int](values)
				data, err := view.MarshalJSON()
				assert.Nil(t, data)
				var unsupported *json.UnsupportedTypeError
				assert.ErrorAs(t, err, &unsupported)
			})
		}
	})
}

func TestReadonly_TraversalOrder(t *testing.T) {
	for source, construct := range readonlySources() {
		t.Run(source, func(t *testing.T) {
			values := construct([]int{1, 2, 3})
			_, _ = values.Shift()
			values.Append(4)
			values.Prepend(0)
			view := deque.AsReadonly[int](values)
			assert.Equal(t, "[0 2 3 4]", view.String())
			data, err := view.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, "[0,2,3,4]", string(data))
		})
	}
}

func readonlySources() map[string]func([]int) deque.Deque[int] {
	return map[string]func([]int) deque.Deque[int]{
		"Array":  func(values []int) deque.Deque[int] { return deque.ArrayFrom(values) },
		"Linked": func(values []int) deque.Deque[int] { return deque.LinkedFrom(values) },
	}
}

func assertDequeBehavior(t *testing.T, d deque.Deque[int]) {
	t.Helper()

	assert.True(t, d.IsEmpty())

	d.Appends(2, 3)
	d.Prepends(0, 1)

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
