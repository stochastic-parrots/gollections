package set

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

type record struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
}

func recordID(value record) int {
	return value.ID
}

func TestNewKeyedHashSet(t *testing.T) {
	set := NewKeyedHashSet(4, recordID)

	assert.NotNil(t, set.values)
	assert.NotNil(t, set.keyOf)
	assert.Empty(t, set.values)
}

func TestNewKeyedHashSetFromSlice(t *testing.T) {
	first := record{ID: 1, Name: "first"}
	duplicate := record{ID: 1, Name: "duplicate"}
	second := record{ID: 2, Name: "second"}

	set := NewKeyedHashSetFromSlice([]record{first, duplicate, second}, recordID)

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []record{first, second}, slices.Collect(set.All()))
}

func TestNewKeyedHashSetFromSeq(t *testing.T) {
	first := record{ID: 1, Name: "first"}
	duplicate := record{ID: 1, Name: "duplicate"}
	second := record{ID: 2, Name: "second"}

	set := NewKeyedHashSetFromSeq(slices.Values([]record{first, duplicate, second}), recordID)

	assert.Equal(t, 2, set.Length())
	assert.ElementsMatch(t, []record{first, second}, slices.Collect(set.All()))
}

func TestKeyedHashSet_Length(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.Equal(t, 2, set.Length())
}

func TestKeyedHashSet_IsEmpty(t *testing.T) {
	set := NewKeyedHashSet(0, recordID)
	assert.True(t, set.IsEmpty())

	set.Add(record{ID: 1})
	assert.False(t, set.IsEmpty())
}

func TestKeyedHashSet_Contains(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1, Name: "stored"}}, recordID)

	assert.True(t, set.Contains(record{ID: 1, Name: "query"}))
	assert.False(t, set.Contains(record{ID: 2}))
}

func TestKeyedHashSet_Add(t *testing.T) {
	set := KeyedHashSet[record, int]{keyOf: recordID}
	first := record{ID: 1, Name: "first"}

	assert.True(t, set.Add(first))
	assert.False(t, set.Add(record{ID: 1, Name: "duplicate"}))

	assert.Equal(t, 1, set.Length())
	assert.Equal(t, []record{first}, slices.Collect(set.All()))
}

func TestKeyedHashSet_Adds(t *testing.T) {
	first := record{ID: 1, Name: "first"}
	set := NewKeyedHashSet(0, recordID)

	assert.Equal(t, 2, set.Adds(first, record{ID: 1, Name: "duplicate"}, record{ID: 2, Name: "second"}))
	assert.Zero(t, set.Adds())
	assert.Zero(t, set.Adds(record{ID: 1}, record{ID: 2}))

	assert.Equal(t, 2, set.Length())
	assert.Contains(t, slices.Collect(set.All()), first)
}

func TestKeyedHashSet_Remove(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.Remove(record{ID: 1, Name: "query"}))
	assert.False(t, set.Remove(record{ID: 1}))
	assert.False(t, set.Contains(record{ID: 1}))
	assert.Equal(t, 1, set.Length())
}

func TestKeyedHashSet_Removes(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}}, recordID)

	assert.Equal(t, 2, set.Removes(record{ID: 2}, record{ID: 2}, record{ID: 4}, record{ID: 3}))
	assert.Zero(t, set.Removes())
	assert.Zero(t, set.Removes(record{ID: 2}, record{ID: 4}))
	assert.Equal(t, []record{{ID: 1}}, slices.Collect(set.All()))
}

func TestKeyedHashSet_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		values := []record{{ID: 1}, {ID: 2}, {ID: 3}}
		set := NewKeyedHashSetFromSlice(values, recordID)

		assert.ElementsMatch(t, values, slices.Collect(set.All()))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)
		count := 0

		set.All()(func(record) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewKeyedHashSet(0, recordID)

		assert.Empty(t, slices.Collect(set.All()))
	})
}

func TestKeyedHashSet_Enumerate(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		input := []record{{ID: 1}, {ID: 2}, {ID: 3}}
		set := NewKeyedHashSetFromSlice(input, recordID)
		var indexes []int
		var values []record

		for idx, value := range set.Enumerate() {
			indexes = append(indexes, idx)
			values = append(values, value)
		}

		assert.Equal(t, []int{0, 1, 2}, indexes)
		assert.ElementsMatch(t, input, values)
	})

	t.Run("PartialIteration", func(t *testing.T) {
		set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)
		count := 0

		set.Enumerate()(func(int, record) bool {
			count++
			return false
		})

		assert.Equal(t, 1, count)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewKeyedHashSet(0, recordID)
		count := 0

		set.Enumerate()(func(int, record) bool {
			count++
			return true
		})

		assert.Zero(t, count)
	})
}

func TestKeyedHashSet_Clear(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)
	storage := set.values

	set.Clear()

	assert.True(t, set.IsEmpty())
	assert.Empty(t, storage)

	value := record{ID: 3}
	set.Add(value)
	assert.Equal(t, value, storage[3])
}

func TestKeyedHashSet_MarshalJSON(t *testing.T) {
	t.Run("Values", func(t *testing.T) {
		values := []record{{ID: 1}, {ID: 2}}
		set := NewKeyedHashSetFromSlice(values, recordID)

		data, err := set.MarshalJSON()
		assert.NoError(t, err)

		var decoded []record
		assert.NoError(t, json.Unmarshal(data, &decoded))
		assert.ElementsMatch(t, values, decoded)
	})

	t.Run("EmptySet", func(t *testing.T) {
		set := NewKeyedHashSet(0, recordID)

		data, err := set.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(data))
	})

	t.Run("UnsupportedValue", func(t *testing.T) {
		set := NewKeyedHashSet(0, func(chan int) int { return 0 })
		set.Add(make(chan int))

		_, err := set.MarshalJSON()
		assert.Error(t, err)
	})
}
