package set_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/set"
	"github.com/stretchr/testify/assert"
)

// Embedding the structural interface hides the concrete value's optional methods.
type structuralSet[T any] struct {
	set.Set[T]
}

var _ set.Set[int] = (*structuralSet[int])(nil)
var _ set.Readonly[int] = (*structuralSet[int])(nil)

var _ json.Marshaler = (*set.HashSet[int])(nil)

var _ json.Marshaler = (*set.KeyedHashSet[int, int])(nil)

type optionalSet[T any] struct {
	set.Set[T]
	data []byte
	err  error
}

func (values optionalSet[T]) MarshalJSON() ([]byte, error) {
	return values.data, values.err
}

type member struct {
	ID     int
	Name   string
	Labels []string
}

func memberID(value member) int {
	return value.ID
}

func TestFactoriesImplementSet(t *testing.T) {
	var _ *set.HashSet[int] = set.NewHashSet[int](0)
	var _ *set.HashSet[int] = set.HashSetFrom[int]([]int{1})
	var _ *set.HashSet[int] = set.HashSetFromSeq[int](slices.Values([]int{1}))
	var _ *set.KeyedHashSet[member, int] = set.NewKeyedHashSet(memberID, 0)
	var _ *set.KeyedHashSet[member, int] = set.KeyedHashSetFrom(memberID, []member{{ID: 1}})
	var _ *set.KeyedHashSet[member, int] = set.KeyedHashSetFromSeq(memberID, slices.Values([]member{{ID: 1}}))
	var _ set.Set[int] = set.NewHashSet[int](0)
	var _ set.Set[member] = set.NewKeyedHashSet(memberID, 0)
}

func TestConcreteZeroValues(t *testing.T) {
	var values set.HashSet[int]

	assert.True(t, values.Add(1))
	assert.Equal(t, 1, values.Adds(2, 1))

	assert.Equal(t, 2, values.Len())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestKeyedHashSetConstructors_NilIdentityFunction(t *testing.T) {
	constructors := map[string]func(){
		"New":     func() { set.NewKeyedHashSet[int, int](nil, 0) },
		"From":    func() { set.KeyedHashSetFrom[int, int](nil, nil) },
		"FromSeq": func() { set.KeyedHashSetFromSeq[int, int](nil, slices.Values([]int(nil))) },
	}
	for name, construct := range constructors {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithValue(t, "set: nil identity function", construct)
		})
	}
}

func TestNewHashSet(t *testing.T) {
	assertSetBehavior(t, set.NewHashSet[int](2))
}

func TestHashSetFrom(t *testing.T) {
	data := []int{1, 2, 1}
	values := set.HashSetFrom[int](data)
	data[0] = 3

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestHashSetFromSeq(t *testing.T) {
	values := set.HashSetFromSeq[int](slices.Values([]int{1, 2, 1}))

	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))
}

func TestNewKeyedHashSet(t *testing.T) {
	values := set.NewKeyedHashSet(memberID, 2)
	first := member{ID: 1, Name: "first", Labels: []string{"a"}}

	assert.True(t, values.Add(first))
	assert.False(t, values.Add(member{ID: 1, Name: "duplicate"}))

	assert.Equal(t, 1, values.Len())
	assert.Equal(t, []member{first}, slices.Collect(values.All()))
}

func TestKeyedHashSetFrom(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	data := []member{first, {ID: 1, Name: "duplicate"}, {ID: 2, Name: "second"}}
	values := set.KeyedHashSetFrom(memberID, data)
	data[0] = member{ID: 3, Name: "changed"}

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestKeyedHashSetFromSeq(t *testing.T) {
	first := member{ID: 1, Name: "first"}
	values := set.KeyedHashSetFromSeq(memberID, slices.Values([]member{
		first,
		{ID: 1, Name: "duplicate"},
		{ID: 2, Name: "second"},
	}))

	assert.ElementsMatch(t, []member{first, {ID: 2, Name: "second"}}, slices.Collect(values.All()))
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() set.Set[int]{
		"HashSetFromNil":        func() set.Set[int] { return set.HashSetFrom([]int(nil)) },
		"HashSetFromEmpty":      func() set.Set[int] { return set.HashSetFrom(make([]int, 0, 4)) },
		"KeyedHashSetFromNil":   func() set.Set[int] { return set.KeyedHashSetFrom(func(x int) int { return x }, []int(nil)) },
		"KeyedHashSetFromEmpty": func() set.Set[int] { return set.KeyedHashSetFrom(func(x int) int { return x }, make([]int, 0, 4)) },
	} {
		t.Run(name, func(t *testing.T) {
			values := construct()
			assert.True(t, values.IsEmpty())
			assert.True(t, values.Add(1))
			assert.True(t, values.Contains(1))
		})
	}
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, set.AsReadonly[int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := set.HashSetFrom[int]([]int{1, 2})
		view := set.AsReadonly[int](mutable)

		_, mutableView := any(view).(set.Set[int])
		_, unmarshals := any(view).(json.Unmarshaler)
		assert.False(t, mutableView)
		assert.False(t, unmarshals)
		assert.True(t, view.Contains(1))
		assert.False(t, view.IsEmpty())
		assert.Equal(t, 2, view.Len())
		assert.ElementsMatch(t, []int{1, 2}, slices.Collect(view.All()))

		var indexes []int
		for idx := range view.Enumerate() {
			indexes = append(indexes, idx)
		}
		assert.Equal(t, []int{0, 1}, indexes)

		data, err := view.MarshalJSON()
		assert.NoError(t, err)
		var decoded []int
		assert.NoError(t, json.Unmarshal(data, &decoded))
		assert.ElementsMatch(t, []int{1, 2}, decoded)

		mutable.Add(3)
		assert.True(t, view.Contains(3))
	})
}

func TestStructuralSet(t *testing.T) {
	values := &structuralSet[int]{Set: set.HashSetFrom([]int{1, 2})}
	assert.NotImplements(t, (*fmt.Stringer)(nil), values)
	assert.NotImplements(t, (*json.Marshaler)(nil), values)
	assert.NotImplements(t, (*json.Unmarshaler)(nil), values)

	view := set.AsReadonly[int](values)
	assert.NotImplements(t, (*set.Set[int])(nil), view)
	assert.NotImplements(t, (*json.Unmarshaler)(nil), view)
	values.Add(3)
	assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(view.All()))
	data, err := json.Marshal(view)
	assert.NoError(t, err)
	var decoded []int
	assert.NoError(t, json.Unmarshal(data, &decoded))
	assert.ElementsMatch(t, []int{1, 2, 3}, decoded)
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
					view := set.AsReadonly[int](values)
					data, err := view.MarshalJSON()
					assert.NoError(t, err)
					if len(test.values) == 0 {
						assert.Equal(t, "[]", string(data))
					} else {
						var decoded []int
						assert.NoError(t, json.Unmarshal(data, &decoded))
						assert.ElementsMatch(t, test.values, decoded)
					}
				})
			}
		})
	}

	t.Run("Custom", func(t *testing.T) {
		want := []byte(" { \"custom\": true }\n")
		values := optionalSet[int]{Set: set.HashSetFrom([]int{1, 2}), data: want}
		view := set.AsReadonly[int](values)
		data, err := view.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, want, data)
	})

	t.Run("CustomError", func(t *testing.T) {
		want := errors.New("custom serialization error")
		values := optionalSet[int]{Set: set.HashSetFrom([]int{1, 2}), err: want}
		view := set.AsReadonly[int](values)
		data, err := view.MarshalJSON()
		assert.Nil(t, data)
		assert.Same(t, want, err)
	})

	t.Run("ElementError", func(t *testing.T) {
		for name, values := range map[string]set.Set[chan int]{
			"Concrete":   set.HashSetFrom([]chan int{make(chan int)}),
			"Structural": &structuralSet[chan int]{Set: set.HashSetFrom([]chan int{make(chan int)})},
		} {
			t.Run(name, func(t *testing.T) {
				view := set.AsReadonly[chan int](values)
				data, err := view.MarshalJSON()
				assert.Nil(t, data)
				var unsupported *json.UnsupportedTypeError
				assert.ErrorAs(t, err, &unsupported)
			})
		}
	})
}

func readonlySources() map[string]func([]int) set.Set[int] {
	return map[string]func([]int) set.Set[int]{
		"Hash": func(values []int) set.Set[int] { return set.HashSetFrom(values) },
		"Keyed": func(values []int) set.Set[int] {
			return set.KeyedHashSetFrom(func(value int) int { return value }, values)
		},
		"Structural": func(values []int) set.Set[int] { return &structuralSet[int]{Set: set.HashSetFrom(values)} },
	}
}

func assertSetBehavior(t *testing.T, values set.Set[int]) {
	t.Helper()

	assert.True(t, values.IsEmpty())
	assert.True(t, values.Add(1))
	assert.Equal(t, 1, values.Adds(2, 1))

	assert.False(t, values.IsEmpty())
	assert.Equal(t, 2, values.Len())
	assert.True(t, values.Contains(1))
	assert.False(t, values.Contains(3))
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(values.All()))

	assert.True(t, values.Remove(1))
	assert.False(t, values.Remove(1))
	assert.Equal(t, 2, values.Adds(3, 4))
	assert.Equal(t, 2, values.Removes(2, 2, 4, 5))

	data, err := values.(json.Marshaler).MarshalJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, `[3]`, string(data))

	values.Clear()
	assert.True(t, values.IsEmpty())

	assert.True(t, values.Add(3))
	assert.True(t, values.Contains(3))
}
