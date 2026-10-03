package list_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stretchr/testify/assert"
)

// Embedding the structural interface hides the concrete value's optional methods.
type structuralList[T any] struct {
	list.List[T]
}

var _ list.List[int] = (*structuralList[int])(nil)
var _ list.Readonly[int] = (*structuralList[int])(nil)

var _ json.Marshaler = (*list.ArrayList[int])(nil)
var _ fmt.Stringer = (*list.ArrayList[int])(nil)
var _ fmt.Formatter = (*list.ArrayList[int])(nil)
var _ json.Unmarshaler = (*list.ArrayList[int])(nil)

var _ json.Marshaler = (*list.LinkedList[int])(nil)
var _ fmt.Stringer = (*list.LinkedList[int])(nil)
var _ fmt.Formatter = (*list.LinkedList[int])(nil)
var _ json.Unmarshaler = (*list.LinkedList[int])(nil)

type optionalList[T any] struct {
	list.List[T]
	data []byte
	err  error
}

func (values optionalList[T]) MarshalJSON() ([]byte, error) {
	return values.data, values.err
}

func (values optionalList[T]) String() string { return "custom formatting" }

func TestConstructorsImplementList(t *testing.T) {
	var _ *list.ArrayList[int] = list.NewArray[int](0)
	var _ *list.ArrayList[int] = list.ArrayFrom[int]([]int{1})
	var _ *list.ArrayList[int] = list.ArrayClone[int]([]int{1})
	var _ *list.ArrayList[int] = list.ArrayFromSeq[int](slices.Values([]int{1}))
	var _ *list.LinkedList[int] = list.NewLinked[int]()
	var _ *list.LinkedList[int] = list.LinkedFrom[int]([]int{1})
	var _ *list.LinkedList[int] = list.LinkedFromSeq[int](slices.Values([]int{1}))
	var _ list.List[int] = list.NewArray[int](0)
	var _ list.List[int] = list.NewLinked[int]()
}

func TestConcreteZeroValues(t *testing.T) {
	var array list.ArrayList[int]
	array.Appends(1, 2)
	assert.Equal(t, []int{1, 2}, array.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &array))
	assert.Equal(t, []int{3, 4}, array.ToSlice())

	var linked list.LinkedList[int]
	linked.Appends(1, 2)
	linked.Reverse()
	assert.Equal(t, []int{2, 1}, linked.ToSlice())
	assert.NoError(t, json.Unmarshal([]byte(`[3,4]`), &linked))
	assert.Equal(t, []int{3, 4}, linked.ToSlice())
}

func TestArrayFrom(t *testing.T) {
	data := []int{1, 2}
	list := list.ArrayFrom[int](data)

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestArrayClone(t *testing.T) {
	data := []int{1, 2}
	list := list.ArrayClone[int](data)

	assert.NoError(t, list.Set(0, 10))
	assert.Equal(t, []int{1, 2}, data)
	data[1] = 20
	assert.Equal(t, []int{10, 2}, list.ToSlice())
}

func TestArrayFromSeq(t *testing.T) {
	list := list.ArrayFromSeq[int](slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestLinkedFrom(t *testing.T) {
	data := []int{1, 2}
	list := list.LinkedFrom[int](data)
	data[0] = 10

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestLinkedFromSeq(t *testing.T) {
	list := list.LinkedFromSeq[int](slices.Values([]int{1, 2}))

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestNewArray(t *testing.T) {
	assertListBehavior(t, list.NewArray[int](1))
}

func TestNewLinked(t *testing.T) {
	assertListBehavior(t, list.NewLinked[int]())
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() list.List[int]{
		"ArrayFromNil":    func() list.List[int] { return list.ArrayFrom([]int(nil)) },
		"ArrayFromEmpty":  func() list.List[int] { return list.ArrayFrom(make([]int, 0, 4)) },
		"ArrayCloneNil":   func() list.List[int] { return list.ArrayClone([]int(nil)) },
		"ArrayCloneEmpty": func() list.List[int] { return list.ArrayClone(make([]int, 0, 4)) },
		"LinkedFromNil":   func() list.List[int] { return list.LinkedFrom([]int(nil)) },
		"LinkedFromEmpty": func() list.List[int] { return list.LinkedFrom(make([]int, 0, 4)) },
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
		assert.Nil(t, list.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := list.NewArray[int](0)
		mutable.Appends(1, 2)

		view := list.AsReadonly[int](mutable)
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, unmarshals)

		assert.Equal(t, []int{1, 2}, slices.Collect(view.All()))
		assert.Equal(t, []int{2, 1}, slices.Collect(view.Backward()))
		assert.Equal(t, []int{0, 1}, collectIndexes(view.Enumerate()))
		assert.Equal(t, 2, view.Len())
		assert.False(t, view.IsEmpty())

		x, err := view.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, x)

		idx, ok := view.Find(func(x int) bool { return x == 2 })
		assert.True(t, ok)
		assert.Equal(t, 1, idx)

		assert.True(t, view.Contains(func(x int) bool { return x == 2 }))
		assert.Equal(t, "[1 2]", view.String())

		mutable.Append(3)
		assert.Equal(t, []int{1, 2, 3}, view.ToSlice())

		data, err := json.Marshal(view)
		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})
}

func TestList_Find(t *testing.T) {
	for _, test := range []struct {
		name string
		list list.List[int]
	}{
		{"Array", list.NewArray[int](0)},
		{"Linked", list.NewLinked[int]()},
	} {
		t.Run(test.name, func(t *testing.T) {
			match := func(value int) bool { return value%2 == 0 }
			idx, ok := test.list.Find(match)
			assert.Equal(t, -1, idx)
			assert.False(t, ok)

			test.list.Appends(1, 2, 4, 6)
			idx, ok = test.list.Find(match)
			assert.Equal(t, 1, idx)
			assert.True(t, ok)

			test.list.Reverse()
			idx, ok = test.list.Find(match)
			assert.Equal(t, 0, idx)
			assert.True(t, ok)

			idx, ok = test.list.Find(func(value int) bool { return value < 0 })
			assert.Equal(t, -1, idx)
			assert.False(t, ok)
		})
	}
}

func TestList_Contains(t *testing.T) {
	for _, test := range []struct {
		name string
		list list.List[int]
	}{
		{"Array", list.NewArray[int](0)},
		{"Linked", list.NewLinked[int]()},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			match := func(value int) bool {
				calls++
				return value%2 == 0
			}
			assert.False(t, test.list.Contains(match))
			assert.Equal(t, 0, calls)

			test.list.Appends(1, 2, 4)
			assert.True(t, test.list.Contains(match))
			assert.Equal(t, 2, calls)

			test.list.Reverse()
			calls = 0
			assert.True(t, test.list.Contains(match))
			assert.Equal(t, 1, calls)
			assert.False(t, test.list.Contains(func(value int) bool { return value < 0 }))
		})
	}
}

func TestStructuralList(t *testing.T) {
	values := &structuralList[int]{List: list.ArrayFrom([]int{1, 2})}
	assert.NotImplements(t, (*fmt.Stringer)(nil), values)
	assert.NotImplements(t, (*json.Marshaler)(nil), values)
	assert.NotImplements(t, (*json.Unmarshaler)(nil), values)

	view := list.AsReadonly[int](values)
	assert.NotImplements(t, (*list.List[int])(nil), view)
	assert.NotImplements(t, (*json.Unmarshaler)(nil), view)
	values.Append(3)
	assert.Equal(t, []int{1, 2, 3}, slices.Collect(view.All()))
	data, err := json.Marshal(view)
	assert.NoError(t, err)
	assert.Equal(t, "[1,2,3]", string(data))
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
					view := list.AsReadonly[int](values)
					assert.Equal(t, test.want, view.String())
					assert.Equal(t, test.want, fmt.Sprint(view))
				})
			}
		})
	}

	t.Run("Custom", func(t *testing.T) {
		values := optionalList[int]{List: list.ArrayFrom([]int{1, 2})}
		view := list.AsReadonly[int](values)
		assert.Equal(t, "custom formatting", view.String())
	})
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
					view := list.AsReadonly[int](values)
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

	t.Run("Custom", func(t *testing.T) {
		want := []byte(" { \"custom\": true }\n")
		values := optionalList[int]{List: list.ArrayFrom([]int{1, 2}), data: want}
		view := list.AsReadonly[int](values)
		data, err := view.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, want, data)
	})

	t.Run("CustomError", func(t *testing.T) {
		want := errors.New("custom serialization error")
		values := optionalList[int]{List: list.ArrayFrom([]int{1, 2}), err: want}
		view := list.AsReadonly[int](values)
		data, err := view.MarshalJSON()
		assert.Nil(t, data)
		assert.Same(t, want, err)
	})

	t.Run("ElementError", func(t *testing.T) {
		for name, values := range map[string]list.List[chan int]{
			"Concrete":   list.ArrayFrom([]chan int{make(chan int)}),
			"Structural": &structuralList[chan int]{List: list.ArrayFrom([]chan int{make(chan int)})},
		} {
			t.Run(name, func(t *testing.T) {
				view := list.AsReadonly[chan int](values)
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
			values.Reverse()
			view := list.AsReadonly[int](values)
			assert.Equal(t, "[3 2 1]", view.String())
			data, err := view.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, "[3,2,1]", string(data))
		})
	}
}

func readonlySources() map[string]func([]int) list.List[int] {
	return map[string]func([]int) list.List[int]{
		"Array":      func(values []int) list.List[int] { return list.ArrayFrom(values) },
		"Linked":     func(values []int) list.List[int] { return list.LinkedFrom(values) },
		"Structural": func(values []int) list.List[int] { return &structuralList[int]{List: list.ArrayFrom(values)} },
	}
}

func assertListBehavior(t *testing.T, l list.List[int]) {
	t.Helper()

	assert.True(t, l.IsEmpty())

	l.Appends(1, 3)
	err := l.Insert(1, 2)
	assert.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3}, l.ToSlice())
	assert.Equal(t, []int{1, 2, 3}, slices.Collect(l.All()))
	assert.Equal(t, []int{3, 2, 1}, slices.Collect(l.Backward()))

	x, err := l.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, x)

	idx, ok := l.Find(func(x int) bool { return x == 3 })
	assert.True(t, ok)
	assert.Equal(t, 2, idx)

	assert.True(t, l.Contains(func(x int) bool { return x == 2 }))

	err = l.Set(1, 20)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 20, 3}, l.ToSlice())

	l.Reverse()
	assert.Equal(t, []int{3, 20, 1}, l.ToSlice())

	removed, err := l.Remove(1)
	assert.NoError(t, err)
	assert.Equal(t, 20, removed)
	assert.Equal(t, []int{3, 1}, l.ToSlice())

	data, err := json.Marshal(l)
	assert.NoError(t, err)
	assert.JSONEq(t, `[3,1]`, string(data))

	err = json.Unmarshal([]byte(`[8,9]`), l)
	assert.NoError(t, err)
	assert.Equal(t, []int{8, 9}, l.ToSlice())

	l.Clear()
	assert.True(t, l.IsEmpty())
	assert.Nil(t, l.ToSlice())
}

func collectIndexes(seq func(func(int, int) bool)) []int {
	var indexes []int
	for idx := range seq {
		indexes = append(indexes, idx)
	}
	return indexes
}
