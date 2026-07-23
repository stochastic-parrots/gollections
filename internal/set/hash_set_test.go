package set

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashSet(t *testing.T) {
	set := NewHashSet[int](4)

	assert.NotNil(t, set.values)
	assert.Empty(t, set.values)
}

func TestNewHashSetFromSlice(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 1})

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestNewHashSetFromSeq(t *testing.T) {
	set := NewHashSetFromSeq(slices.Values([]int{1, 2, 1}))

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestHashSet_Length(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.Equal(t, 2, set.Length())
}

func TestHashSet_IsEmpty(t *testing.T) {
	set := NewHashSet[int](0)
	assert.True(t, set.IsEmpty())

	set.Add(1)
	assert.False(t, set.IsEmpty())
}

func TestHashSet_Contains(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.Contains(1))
	assert.False(t, set.Contains(3))
}

func TestHashSet_Add(t *testing.T) {
	var set HashSet[int]

	assert.True(t, set.Add(1))
	assert.False(t, set.Add(1))

	assert.Equal(t, 1, set.Length())
	assert.True(t, set.Contains(1))
}

func TestHashSet_Adds(t *testing.T) {
	set := NewHashSet[int](0)

	assert.Equal(t, 2, set.Adds(1, 2, 1))
	assert.Zero(t, set.Adds())
	assert.Zero(t, set.Adds(1, 2))

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []int{1, 2}, slices.Collect(set.All()))
}

func TestHashSet_Remove(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2})

	assert.True(t, set.Remove(1))
	assert.False(t, set.Remove(1))
	assert.False(t, set.Contains(1))
	assert.Equal(t, 1, set.Length())
}

func TestHashSet_Removes(t *testing.T) {
	set := NewHashSetFromSlice([]int{1, 2, 3})

	assert.Equal(t, 2, set.Removes(2, 2, 4, 3))
	assert.Zero(t, set.Removes())
	assert.Zero(t, set.Removes(2, 4))
	assert.Equal(t, []int{1}, slices.Collect(set.All()))
}

func TestHashSet_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2, 3})

		assert.ElementsMatch(t, []int{1, 2, 3}, slices.Collect(set.All()))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2, 3})
		count := 0

		set.All()(func(int) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		assert.Empty(t, slices.Collect(set.All()))
	})
}

func TestHashSet_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{10, 20, 30})
		var indexes []int
		var values []int

		for idx, value := range set.Enumerate() {
			indexes = append(indexes, idx)
			values = append(values, value)
		}

		assert.Equal(t, []int{0, 1, 2}, indexes)
		assert.ElementsMatch(t, []int{10, 20, 30}, values)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{10, 20, 30})
		count := 0

		set.Enumerate()(func(int, int) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		assert.Empty(t, slices.Collect(func(yield func(int) bool) {
			for _, value := range set.Enumerate() {
				if !yield(value) {
					return
				}
			}
		}))
	})
}

func TestHashSet_Clear(t *testing.T) {
	first, second, third := 1, 2, 3
	set := NewHashSetFromSlice([]*int{&first, &second})
	storage := set.values

	set.Clear()

	assert.True(t, set.IsEmpty())
	assert.Empty(t, storage)

	set.Add(&third)
	_, reused := storage[&third]
	assert.True(t, reused)
}

func TestHashSet_MarshalJSON(t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		set := NewHashSetFromSlice([]int{1, 2})

		data, err := set.MarshalJSON()
		assert.NoError(t, err)

		var values []int
		assert.NoError(t, json.Unmarshal(data, &values))
		assert.ElementsMatch(t, []int{1, 2}, values)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewHashSet[int](0)

		data, err := set.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("UnsupportedValue", func(t *testing.T) {
		set := NewHashSet[chan int](0)
		set.Add(make(chan int))

		_, err := set.MarshalJSON()
		assert.Error(t, err)
	})
}
