package prioritymap_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/prioritymap"
	"github.com/stretchr/testify/assert"
)

func TestFactoriesImplementPriorityMap(t *testing.T) {
	var _ prioritymap.PriorityMap[string, int] = prioritymap.NewBinaryHeap[string](0, cmp.Less[int])
	var _ prioritymap.PriorityMap[string, int] = prioritymap.MinBinaryHeap[string, int](0)
	var _ prioritymap.PriorityMap[string, int] = prioritymap.MaxBinaryHeap[string, int](0)
	var _ prioritymap.PriorityMap[string, int] = prioritymap.NewPairingHeap[string](0, cmp.Less[int])
	var _ prioritymap.PriorityMap[string, int] = prioritymap.MinPairingHeap[string, int](0)
	var _ prioritymap.PriorityMap[string, int] = prioritymap.MaxPairingHeap[string, int](0)
	var _ prioritymap.PriorityMap[string, uint64] = prioritymap.NewRadixHeap[string, uint64](0)
}

func TestNewBinaryHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.NewBinaryHeap[string](0, cmp.Less[int]), "one", 1, 0, []int{0, 1, 3})
}

func TestMinBinaryHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.MinBinaryHeap[string, int](0), "one", 1, 0, []int{0, 1, 3})
}

func TestMaxBinaryHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.MaxBinaryHeap[string, int](0), "three", 3, 4, []int{4, 3, 1})
}

func TestNewPairingHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.NewPairingHeap[string](0, cmp.Less[int]), "one", 1, 0, []int{0, 1, 3})
}

func TestMinPairingHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.MinPairingHeap[string, int](0), "one", 1, 0, []int{0, 1, 3})
}

func TestMaxPairingHeap(t *testing.T) {
	assertPriorityMapBehavior(t, prioritymap.MaxPairingHeap[string, int](0), "three", 3, 4, []int{4, 3, 1})
}

func TestNewRadixHeap(t *testing.T) {
	type distance uint32
	pm := prioritymap.NewRadixHeap[string, distance](0)
	pm.Set("far", 100)
	pm.Set("near", 10)
	pm.Set("middle", 50)

	key, value, ok := pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, "near", key)
	assert.Equal(t, distance(10), value)

	assert.True(t, pm.Improve("far", 60))
	assert.False(t, pm.Improve("far", 70))

	var priorities []distance
	for _, priority := range pm.Drain() {
		priorities = append(priorities, priority)
	}
	assert.Equal(t, []distance{50, 60}, priorities)
}

func TestAsReadonly(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, prioritymap.AsReadonly[string, int](nil))
	})

	t.Run("View", func(t *testing.T) {
		mutable := prioritymap.MinBinaryHeap[string, int](0)
		mutable.Set("slow", 10)
		mutable.Set("fast", 1)

		view := prioritymap.AsReadonly(mutable)

		assert.Equal(t, 2, view.Length())
		assert.True(t, view.Contains("slow"))

		priority, ok := view.Get("slow")
		assert.True(t, ok)
		assert.Equal(t, 10, priority)

		key, priority, ok := view.Peek()
		assert.True(t, ok)
		assert.Equal(t, "fast", key)
		assert.Equal(t, 1, priority)

		mutable.Set("urgent", 0)
		key, priority, ok = view.Peek()
		assert.True(t, ok)
		assert.Equal(t, "urgent", key)
		assert.Equal(t, 0, priority)

		assert.Len(t, slices.Collect(view.Keys()), 3)
		assert.Len(t, slices.Collect(view.Values()), 3)
		assert.Equal(t, 3, countAll(view.All()))
	})
}

func assertPriorityMapBehavior(
	t *testing.T,
	pm prioritymap.PriorityMap[string, int],
	initialBestKey string,
	initialBestPriority int,
	improvePriority int,
	expectedDrain []int,
) {
	t.Helper()

	assert.True(t, pm.IsEmpty())

	pm.Set("one", 1)
	pm.Set("two", 2)
	pm.Set("three", 3)

	assert.Equal(t, 3, pm.Length())
	assert.True(t, pm.Contains("two"))
	assert.False(t, pm.Contains("missing"))

	priority, ok := pm.Get("two")
	assert.True(t, ok)
	assert.Equal(t, 2, priority)

	key, priority, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, initialBestKey, key)
	assert.Equal(t, initialBestPriority, priority)

	assert.True(t, pm.Update("three", 3))
	assert.False(t, pm.Update("missing", 0))

	pm.Set("removed", 2)
	assert.True(t, pm.Remove("removed"))
	assert.False(t, pm.Remove("missing"))

	assert.True(t, pm.Improve("two", improvePriority))
	assert.False(t, pm.Improve("two", improvePriority))

	key, priority, ok = pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "two", key)
	assert.Equal(t, improvePriority, priority)

	assert.Len(t, slices.Collect(pm.Keys()), 3)
	assert.Len(t, slices.Collect(pm.Values()), 3)
	assert.Equal(t, 3, countAll(pm.All()))

	assert.Equal(t, expectedDrain, drainPriorities(pm))
	assert.True(t, pm.IsEmpty())

	pm.Set("one", 1)
	pm.Clear()
	assert.True(t, pm.IsEmpty())
}

func countAll(seq func(func(string, int) bool)) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

func drainPriorities(pm prioritymap.PriorityMap[string, int]) []int {
	var priorities []int
	for _, priority := range pm.Drain() {
		priorities = append(priorities, priority)
	}
	return priorities
}
