package prioritymap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stretchr/testify/assert"
)

func TestNewBinaryPriorityMap(t *testing.T) {
	pm := NewBinaryPriorityMap[string](10, comparator.Min[int]())

	assert.Equal(t, 0, pm.Length())
	assert.True(t, pm.IsEmpty())
	assert.Empty(t, pm.data)
	assert.Empty(t, pm.indexes)
}

func TestBinaryPriorityMapSet(t *testing.T) {
	t.Run("Insert New Keys", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, 5, val)
	})

	t.Run("Update Existing Keys", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2) // Update priority to be higher

		assert.Equal(t, 1, pm.Length())
		_, priority, _ := pm.Peek()
		assert.Equal(t, 2, priority)
	})
}

func TestBinaryPriorityMapUpdate(t *testing.T) {
	t.Run("Update nonexistent key", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		assert.False(t, pm.Update("nonexistent", 1))
	})

	t.Run("Update existent key", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("apple", 50)
		pm.Set("banana", 30)
		pm.Set("cherry", 10)

		assert.True(t, pm.Update("apple", 5))

		key, priority, ok := pm.Peek()
		assert.Equal(t, "apple", key)
		assert.Equal(t, 5, priority)
		assert.True(t, ok)
	})
}

func TestBinaryPriorityMapGet(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestBinaryPriorityMapRemove(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("A", 10)
	pm.Set("B", 20)
	pm.Set("C", 30)

	assert.True(t, pm.Remove("B"))
	assert.Equal(t, 2, pm.Length())
	_, exists := pm.Get("B")
	assert.False(t, exists)

	assert.False(t, pm.Remove("non-existent"))
}

func TestBinaryPriorityMapPop(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	pm.Set(1, 50)
	pm.Set(2, 10)
	pm.Set(3, 30)

	// Expected order: 2 (10), 3 (30), 1 (50)
	k, v, ok := pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, 10, v)

	pm.Pop()
	pm.Pop()

	_, _, ok = pm.Pop()
	assert.False(t, ok)
}

func TestBinaryPriorityMapPeek(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())

	_, _, ok := pm.Peek()
	assert.False(t, ok)

	pm.Set("top", 1)
	k, v, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "top", k)
	assert.Equal(t, 1, v)
}

func TestBinaryPriorityMapContains(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("a", 10)

	assert.False(t, !pm.Contains("a"))
	assert.True(t, pm.Contains("a"))
	assert.False(t, pm.Contains("b"))

	pm.Remove("a")
	assert.False(t, pm.Contains("a"))
}

func TestBinaryPriorityMapIsEmpty(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	assert.True(t, pm.IsEmpty())

	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())

	pm.Pop()
	assert.True(t, pm.IsEmpty())
}

func TestBinaryPriorityMapLength(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	assert.Equal(t, 0, pm.Length())

	pm.Set(1, 10)
	pm.Set(2, 20)
	assert.Equal(t, 2, pm.Length())
}

func TestBinaryPriorityMapKeys(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("A", 1)
	pm.Set("B", 2)

	keys := make(map[string]bool)
	for k := range pm.Keys() {
		keys[k] = true
	}

	assert.Len(t, keys, 2)
	assert.True(t, keys["A"])
	assert.True(t, keys["B"])
}

func TestBinaryPriorityMapValues(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("A", 100)
	pm.Set("B", 200)

	var values []int
	for v := range pm.Values() {
		values = append(values, v)
	}

	assert.Len(t, values, 2)
	assert.Contains(t, values, 100)
	assert.Contains(t, values, 200)
}

func TestBinaryPriorityMapAll(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("A", 1)
	pm.Set("B", 2)

	count := 0
	for k, v := range pm.All() {
		val, ok := pm.Get(k)
		assert.True(t, ok)
		assert.Equal(t, val, v)
		count++
	}
	assert.Equal(t, 2, count)
}

func TestBinaryPriorityMapDrain(t *testing.T) {
	t.Run("Total Drain", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Max[int]())
		items := map[string]int{"a": 30, "b": 10, "c": 20}
		for k, v := range items {
			pm.Set(k, v)
		}

		keys := []string{"a", "c", "b"}
		priorities := []int{30, 20, 10}
		idx := 0
		for key, priority := range pm.Drain() {
			assert.Equal(t, keys[idx], key)
			assert.Equal(t, priorities[idx], priority)
			idx++
		}

		assert.True(t, pm.IsEmpty())
	})

	t.Run("Partial Drain (break)", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		for key := range pm.Drain() {
			if key == "a" {
				break
			}
		}

		assert.Equal(t, 2, pm.Length())
		assert.False(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("Drain Empty Map", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		count := 0
		for range pm.Drain() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestBinaryPriorityMapIntegrity(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())

	for i := 10; i > 0; i-- {
		pm.Set(i, i)
	}

	pm.Set(5, 100)
	pm.Remove(1)
	pm.Set(10, 0)

	for key, idx := range pm.indexes {
		assert.Equal(t, key, pm.data[idx].key, "Index map out of sync for key %d", key)
	}
}

func TestBinaryPriorityMapClear(t *testing.T) {
	t.Run("Clear populated map", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](10, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)
		assert.Equal(t, 3, pm.Length())

		pm.Clear()

		assert.Equal(t, 0, pm.Length())
		assert.True(t, pm.IsEmpty())
		assert.False(t, pm.Contains("a"))
		assert.False(t, pm.Contains("b"))
		assert.False(t, pm.Contains("c"))

		_, _, ok := pm.Peek()
		assert.False(t, ok)

		_, _, ok = pm.Pop()
		assert.False(t, ok)

		assert.Equal(t, 0, len(pm.indexes))
		assert.Equal(t, 0, len(pm.data))
		assert.Equal(t, 10, cap(pm.data))
	})

	t.Run("Reuse after Clear", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())

		pm.Set("old", 100)
		pm.Clear()
		assert.Equal(t, 0, len(pm.indexes))
		assert.Equal(t, 0, len(pm.data))
		assert.Equal(t, 1, cap(pm.data))

		pm.Set("new", 5)
		pm.Set("newer", 1)

		assert.Equal(t, 2, pm.Length())
		key, priority, ok := pm.Peek()
		assert.True(t, ok)
		assert.Equal(t, "newer", key)
		assert.Equal(t, 1, priority)
	})
}
