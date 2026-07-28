package set

import (
	"encoding/json"
	"iter"
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

type keyedSource[T any] []T

func (values keyedSource[T]) Length() int {
	return len(values)
}

func (values keyedSource[T]) All() iter.Seq[T] {
	return slices.Values(values)
}

func recordNameLength(value record) int {
	return len(value.Name)
}

func TestCollectKeyedKeys(t *testing.T) {
	differentIdentity := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, recordNameLength)

	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(
		keyedSource[record]{{ID: 1}, {ID: 1}, {ID: 2}},
		recordID,
	))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	))
	assert.Equal(t, map[int]struct{}{
		1: {},
		2: {},
	}, CollectKeyedKeys(differentIdentity, recordID))
}

func TestCollectKeyedMembership(t *testing.T) {
	differentIdentity := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, recordNameLength)

	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(
		keyedSource[record]{{ID: 1}, {ID: 1}, {ID: 2}},
		recordID,
	))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	))
	assert.Equal(t, map[int]bool{
		1: false,
		2: false,
	}, CollectKeyedMembership(differentIdentity, recordID))
}

func TestMatchKeyedMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, contained := MatchKeyedMembership(
		membership,
		keyedSource[record]{{ID: 2}, {ID: 2}, {ID: 1}},
		recordID,
	)

	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false},
		keyedSource[record]{{ID: 1}, {ID: 2}},
		recordID,
	)
	assert.Zero(t, matched)
	assert.False(t, contained)

	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 2}, {ID: 1}}, recordID),
		recordID,
	)
	assert.Equal(t, 2, matched)
	assert.True(t, contained)
	matched, contained = MatchKeyedMembership(
		map[int]bool{1: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
		recordID,
	)
	assert.Zero(t, matched)
	assert.False(t, contained)
}

func TestMatchKeyedSubset(t *testing.T) {
	assert.True(t, MatchKeyedSubset(
		map[int]bool{},
		keyedSource[record]{{ID: 1}},
		recordID,
	))
	assert.True(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		keyedSource[record]{{ID: 2}, {ID: 2}, {ID: 3}, {ID: 1}},
		recordID,
	))
	assert.False(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		keyedSource[record]{{ID: 1}, {ID: 3}},
		recordID,
	))
	assert.True(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice(
			[]record{{ID: 2}, {ID: 3}, {ID: 1}},
			recordID,
		),
		recordID,
	))
	assert.False(t, MatchKeyedSubset(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 3}}, recordID),
		recordID,
	))
}

func TestMatchKeyedContainedMembership(t *testing.T) {
	membership := map[int]bool{1: false, 2: false}

	matched, hasExtra := MatchKeyedContainedMembership(
		membership,
		keyedSource[record]{{ID: 1}, {ID: 3}, {ID: 2}, {ID: 4}},
		recordID,
	)

	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	assert.Equal(t, map[int]bool{1: true, 2: true}, membership)

	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		keyedSource[record]{{ID: 1}},
		recordID,
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)

	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice(
			[]record{{ID: 1}, {ID: 3}, {ID: 2}},
			recordID,
		),
		recordID,
	)
	assert.Equal(t, 2, matched)
	assert.True(t, hasExtra)
	matched, hasExtra = MatchKeyedContainedMembership(
		map[int]bool{1: false, 2: false},
		NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID),
		recordID,
	)
	assert.Equal(t, 1, matched)
	assert.False(t, hasExtra)
}

func TestNewKeyedHashSetIntersection(t *testing.T) {
	result := NewKeyedHashSetIntersection(
		recordID,
		keyedSource[record]{{ID: 1}, {ID: 2, Name: "first"}, {ID: 2}},
		keyedSource[record]{{ID: 2}, {ID: 3}},
		keyedSource[record]{{ID: 2}},
	)

	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(result.All()))
	assert.True(t, NewKeyedHashSetIntersection(recordID).IsEmpty())
	assert.True(t, NewKeyedHashSetIntersection(
		recordID,
		keyedSource[record]{},
		keyedSource[record]{{ID: 1}},
	).IsEmpty())
	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 2}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}},
		).All(),
	))
	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			keyedSource[record]{{ID: 2, Name: "first"}},
			keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 3}},
		).All(),
	))
	filter := NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID)
	assert.Equal(t, []record{{ID: 2, Name: "first"}}, slices.Collect(
		NewKeyedHashSetIntersection(
			recordID,
			keyedSource[record]{{ID: 1}, {ID: 2, Name: "first"}},
			filter,
		).All(),
	))
	assert.True(t, NewKeyedHashSetIntersection(
		recordID,
		keyedSource[record]{{ID: 1}, {ID: 2}},
		keyedSource[record]{{ID: 1}},
		keyedSource[record]{{ID: 2}},
	).IsEmpty())
}

func TestNewKeyedHashSetDifference(t *testing.T) {
	assert.True(t, NewKeyedHashSetDifference(recordID).IsEmpty())
	assert.Equal(t, []record{{ID: 1, Name: "first"}}, slices.Collect(
		NewKeyedHashSetDifference(
			recordID,
			keyedSource[record]{{ID: 1, Name: "first"}, {ID: 2}},
			NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID),
		).All(),
	))
	assert.True(t, NewKeyedHashSetDifference(
		recordID,
		keyedSource[record]{{ID: 1}, {ID: 2}},
		keyedSource[record]{{ID: 1}, {ID: 2}},
		keyedSource[record]{{ID: 3}},
	).IsEmpty())
}

func TestNewKeyedHashSetSymmetricDifference(t *testing.T) {
	result := NewKeyedHashSetSymmetricDifference(
		recordID,
		keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}},
		keyedSource[record]{{ID: 2}, {ID: 3, Name: "three"}},
	)

	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 3, Name: "three"}}, slices.Collect(result.All()))
	assert.True(t, NewKeyedHashSetSymmetricDifference(recordID).IsEmpty())
}

func TestKeyedHashSet_Clone(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
	}, recordID)
	clone := original.Clone()

	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(clone.All()))
	assert.True(t, clone.Add(record{ID: 3, Name: "three"}))
	assert.False(t, original.Contains(record{ID: 3}))
	assert.Equal(t, 3, clone.keyOf(record{ID: 3}))
}

func TestKeyedHashSet_Union(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
	}, recordID)
	result := original.Union(keyedSource[record]{
		{ID: 2, Name: "right"},
		{ID: 3, Name: "three"},
		{ID: 3, Name: "duplicate"},
	})
	self := original.Union(original)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Union().All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(self.All()))
	assert.True(t, self.Add(record{ID: 3}))
	assert.False(t, original.Contains(record{ID: 3}))
}

func TestKeyedHashSet_Intersection(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
	}, recordID)
	result := original.Intersection(
		keyedSource[record]{{ID: 2, Name: "right"}, {ID: 3}},
		keyedSource[record]{{ID: 2, Name: "again"}},
	)

	assert.Equal(t, []record{{ID: 2, Name: "left"}}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Intersection().All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(
		original.Intersection(original).All(),
	))
	assert.True(t, original.Intersection(keyedSource[record]{}).IsEmpty())
	assert.True(t, original.Intersection(
		keyedSource[record]{{ID: 1}},
		keyedSource[record]{{ID: 2}},
	).IsEmpty())
	assert.Equal(t, []record{{ID: 2, Name: "left"}}, slices.Collect(
		original.Intersection(NewKeyedHashSetFromSlice(
			[]record{{ID: 2}},
			recordID,
		)).All(),
	))

	empty := NewKeyedHashSet(0, recordID)
	assert.True(t, empty.Intersection(keyedSource[record]{{ID: 1}}).IsEmpty())
}

func TestKeyedHashSet_Difference(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
		{ID: 3, Name: "three"},
	}, recordID)
	result := original.Difference(
		keyedSource[record]{{ID: 2}, {ID: 2}},
		keyedSource[record]{{ID: 4}},
	)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "three"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.Difference().All()))
	assert.ElementsMatch(t, []record{{ID: 1, Name: "one"}, {ID: 3, Name: "three"}}, slices.Collect(
		original.Difference(NewKeyedHashSetFromSlice([]record{{ID: 2}}, recordID)).All(),
	))
	assert.True(t, original.Difference(original).IsEmpty())
	assert.True(t, original.Difference(
		keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 3}},
	).IsEmpty())
}

func TestKeyedHashSet_SymmetricDifference(t *testing.T) {
	original := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "left"},
	}, recordID)
	result := original.SymmetricDifference(
		keyedSource[record]{{ID: 2}, {ID: 3, Name: "first"}},
		keyedSource[record]{{ID: 3, Name: "second"}},
		keyedSource[record]{{ID: 3, Name: "third"}, {ID: 4, Name: "four"}},
	)

	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "third"},
		{ID: 4, Name: "four"},
	}, slices.Collect(result.All()))
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(original.SymmetricDifference().All()))
	assert.True(t, original.SymmetricDifference(original).IsEmpty())
	assert.ElementsMatch(t, slices.Collect(original.All()), slices.Collect(
		original.SymmetricDifference(original, original).All(),
	))
	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3},
	}, slices.Collect(original.SymmetricDifference(
		NewKeyedHashSetFromSlice([]record{{ID: 2}, {ID: 3}}, recordID),
	).All()))
	assert.Equal(t, []record{
		{ID: 1, Name: "one"},
	}, slices.Collect(original.SymmetricDifference(
		keyedSource[record]{{ID: 2}, {ID: 2}},
	).All()))
	byNameLength := func(value record) int {
		return len(value.Name)
	}
	assert.Equal(t, 1, original.SymmetricDifference(
		NewKeyedHashSetFromSlice([]record{
			{ID: 1, Name: "a"},
			{ID: 1, Name: "bb"},
		}, byNameLength),
	).Length())
}

func TestKeyedHashSet_RelationsReapplyReceiverIdentity(t *testing.T) {
	byID := NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID)
	byNameLength := func(value record) int {
		return len(value.Name)
	}
	sameIDs := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
	}, byNameLength)

	assert.True(t, byID.Equal(sameIDs))
	assert.True(t, byID.IsSubset(sameIDs))
	assert.False(t, byID.IsProperSubset(sameIDs))
	assert.True(t, byID.IsSuperset(sameIDs))
	assert.False(t, byID.IsProperSuperset(sameIDs))
	assert.False(t, byID.IsDisjoint(sameIDs))

	withExtra := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 2, Name: "bb"},
	}, byNameLength)
	assert.False(t, byID.Equal(withExtra))
	assert.True(t, byID.IsProperSubset(withExtra))
	assert.False(t, byID.IsSuperset(withExtra))
	byIDWithExtra := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)
	assert.True(t, byIDWithExtra.IsProperSuperset(sameIDs))
	assert.False(t, byIDWithExtra.IsSubset(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}, {ID: 3}},
		recordID,
	)))
	assert.False(t, byIDWithExtra.IsProperSuperset(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}, {ID: 3}},
		recordID,
	)))
	assert.True(t, byID.IsDisjoint(NewKeyedHashSetFromSlice(
		[]record{{ID: 2, Name: "a"}},
		byNameLength,
	)))
}

func TestKeyedHashSet_Equal(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.Equal(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.True(t, set.Equal(set))
	assert.False(t, set.Equal(keyedSource[record]{{ID: 1}}))
	assert.False(t, set.Equal(keyedSource[record]{{ID: 1}, {ID: 3}}))
}

func TestKeyedHashSet_IsSubset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.IsSubset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 3}}))
	assert.True(t, set.IsSubset(set))
	assert.False(t, set.IsSubset(keyedSource[record]{{ID: 1}}))
	assert.False(t, set.IsSubset(keyedSource[record]{{ID: 1}, {ID: 3}}))
}

func TestKeyedHashSet_IsProperSubset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.False(t, set.IsProperSubset(keyedSource[record]{{ID: 1}}))
	assert.True(t, set.IsProperSubset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 3}}))
	assert.False(t, set.IsProperSubset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.False(t, set.IsProperSubset(keyedSource[record]{{ID: 1}, {ID: 3}, {ID: 4}}))
}

func TestKeyedHashSet_IsSuperset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}}, recordID)

	assert.True(t, set.IsSuperset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.True(t, set.IsSuperset(set))
	assert.False(t, set.IsSuperset(keyedSource[record]{{ID: 1}, {ID: 4}}))
}

func TestKeyedHashSet_IsProperSuperset(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}}, recordID)

	assert.True(t, set.IsProperSuperset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 2}}))
	assert.False(t, set.IsProperSuperset(keyedSource[record]{{ID: 1}, {ID: 2}, {ID: 3}}))
	assert.False(t, set.IsProperSuperset(set))
	assert.False(t, set.IsProperSuperset(keyedSource[record]{{ID: 1}, {ID: 4}}))
}

func TestKeyedHashSet_IsDisjoint(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID)

	assert.True(t, set.IsDisjoint(keyedSource[record]{{ID: 3}, {ID: 4}, {ID: 4}}))
	assert.False(t, set.IsDisjoint(keyedSource[record]{{ID: 2}, {ID: 3}}))
	assert.False(t, set.IsDisjoint(set))
	empty := NewKeyedHashSet(0, recordID)
	assert.True(t, empty.IsDisjoint(empty))
}

func TestKeyedHashSet_UnionWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2, Name: "left"}}, recordID)

	assert.Equal(t, 2, set.UnionWith(
		keyedSource[record]{{ID: 2, Name: "right"}, {ID: 3, Name: "three"}},
		keyedSource[record]{{ID: 4, Name: "four"}},
	))
	assert.Zero(t, set.UnionWith())
	assert.Zero(t, set.UnionWith(set))
	assert.Zero(t, set.UnionWith(set, set))
	assert.Equal(t, 1, set.UnionWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 5, Name: "five"}},
		recordID,
	)))
	assert.ElementsMatch(t, []record{
		{ID: 1},
		{ID: 2, Name: "left"},
		{ID: 3, Name: "three"},
		{ID: 4, Name: "four"},
		{ID: 5, Name: "five"},
	}, slices.Collect(set.All()))
}

func TestKeyedHashSet_IntersectWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{
		{ID: 1},
		{ID: 2},
		{ID: 3, Name: "stored"},
		{ID: 4},
	}, recordID)

	assert.Equal(t, 3, set.IntersectWith(
		keyedSource[record]{{ID: 2}, {ID: 3}, {ID: 4}, {ID: 4}},
		keyedSource[record]{{ID: 3, Name: "query"}},
	))
	assert.Zero(t, set.IntersectWith())
	assert.Zero(t, set.IntersectWith(set))
	assert.Equal(t, []record{{ID: 3, Name: "stored"}}, slices.Collect(set.All()))

	assert.Equal(t, 1, set.IntersectWith(keyedSource[record]{}))
	assert.True(t, set.IsEmpty())

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 2, set.IntersectWith(
		keyedSource[record]{{ID: 1}},
		keyedSource[record]{{ID: 2}},
	))
	assert.True(t, set.IsEmpty())
}

func TestKeyedHashSet_DifferenceWith(t *testing.T) {
	set := NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}}, recordID)

	assert.Equal(t, 2, set.DifferenceWith(
		keyedSource[record]{{ID: 2}, {ID: 2}},
		keyedSource[record]{{ID: 4}, {ID: 5}},
	))
	assert.Zero(t, set.DifferenceWith())
	assert.ElementsMatch(t, []record{{ID: 1}, {ID: 3}}, slices.Collect(set.All()))

	assert.Equal(t, 2, set.DifferenceWith(set))
	assert.True(t, set.IsEmpty())
	assert.Zero(t, set.DifferenceWith(keyedSource[record]{{ID: 1}}))

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 2, set.DifferenceWith(keyedSource[record]{{ID: 1}, {ID: 2}}))
	assert.True(t, set.IsEmpty())

	set.Adds(record{ID: 1}, record{ID: 2})
	assert.Equal(t, 1, set.DifferenceWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 2}},
		recordID,
	)))
	assert.Equal(t, []record{{ID: 1}}, slices.Collect(set.All()))
	assert.Equal(t, 1, set.DifferenceWith(NewKeyedHashSetFromSlice(
		[]record{{ID: 1}},
		recordID,
	)))
	assert.True(t, set.IsEmpty())
}

func TestKeyedHashSet_SymmetricDifferenceWith(t *testing.T) {
	set := KeyedHashSet[record, int]{keyOf: recordID}

	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		NewKeyedHashSetFromSlice([]record{{ID: 4}}, recordID),
	))
	assert.Equal(t, []record{{ID: 4}}, slices.Collect(set.All()))
	set.Clear()

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		keyedSource[record]{{ID: 1, Name: "one"}, {ID: 2}, {ID: 2}},
		keyedSource[record]{{ID: 2}, {ID: 3, Name: "three"}, {ID: 3}},
	))
	assert.ElementsMatch(t, []record{
		{ID: 1, Name: "one"},
		{ID: 3, Name: "three"},
	}, slices.Collect(set.All()))
	assert.Zero(t, set.SymmetricDifferenceWith())

	assert.Equal(t, 2, set.SymmetricDifferenceWith(&set))
	assert.True(t, set.IsEmpty())

	otherKey := func(value record) int {
		return len(value.Name)
	}
	other := NewKeyedHashSetFromSlice([]record{
		{ID: 1, Name: "a"},
		{ID: 1, Name: "bb"},
		{ID: 2, Name: "ccc"},
	}, otherKey)

	single := NewKeyedHashSetFromSlice([]record{{ID: 1}}, recordID)
	assert.Equal(t, 2, single.SymmetricDifferenceWith(other))
	assert.Equal(t, []record{{ID: 2, Name: "ccc"}}, slices.Collect(single.All()))

	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		other,
		keyedSource[record]{{ID: 2}},
		keyedSource[record]{{ID: 3}},
	))
	assert.Equal(t, 2, set.Length())
	assert.True(t, set.Contains(record{ID: 1}))
	assert.True(t, set.Contains(record{ID: 3}))

	set.Clear()
	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		keyedSource[record]{{ID: 1}},
		NewKeyedHashSetFromSlice([]record{{ID: 1}, {ID: 2}}, recordID),
	))
	assert.Equal(t, []record{{ID: 2}}, slices.Collect(set.All()))

	set.Clear()
	assert.Equal(t, 1, set.SymmetricDifferenceWith(
		keyedSource[record]{{ID: 2}},
		other,
	))
	assert.Equal(t, 1, set.Length())
	assert.True(t, set.Contains(record{ID: 1}))

	set.Clear()
	set.Add(record{ID: 1})
	assert.Equal(t, 2, set.SymmetricDifferenceWith(
		keyedSource[record]{{ID: 1}},
		keyedSource[record]{{ID: 2}},
	))
	assert.Equal(t, []record{{ID: 2}}, slices.Collect(set.All()))
}
