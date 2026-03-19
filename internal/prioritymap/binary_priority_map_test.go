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

func TestBinaryPriorityMap_Set(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, 5, val)
	})

	t.Run("Update", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2) // Update priority to be higher

		assert.Equal(t, 1, pm.Length())
		_, priority, _ := pm.Peek()
		assert.Equal(t, 2, priority)
	})
}

func TestBinaryPriorityMap_Update(t *testing.T) {
	t.Run("Nonexistent", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		assert.False(t, pm.Update("nonexistent", 1))
	})

	t.Run("Existent", func(t *testing.T) {
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

func TestBinaryPriorityMap_SetIfBetter(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())

		assert.True(t, pm.SetIfBetter("A", 10))
		assert.Equal(t, 1, pm.Length())

		val, ok := pm.Get("A")
		assert.True(t, ok)
		assert.Equal(t, 10, val)
	})

	t.Run("ImprovePriority", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("Root", 1)
		pm.Set("A", 10)
		pm.Set("B", 20)

		assert.True(t, pm.SetIfBetter("A", 5))

		key, priority, _ := pm.Peek()
		assert.Equal(t, "Root", key)
		assert.Equal(t, 1, priority)

		assert.True(t, pm.SetIfBetter("Root", 0))

		key, priority, _ = pm.Peek()
		assert.Equal(t, "Root", key)
		assert.Equal(t, 0, priority)
		priority, _ = pm.Get("A")
		assert.Equal(t, 5, priority)
	})

	t.Run("WorsenPriority", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)

		assert.False(t, pm.SetIfBetter("A", 15))

		val, _ := pm.Get("A")
		assert.Equal(t, 10, val)
	})

	t.Run("SamePriority", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("A", 10)

		assert.False(t, pm.SetIfBetter("A", 10))
	})
}

func TestBinaryPriorityMap_Get(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestBinaryPriorityMap_Remove(t *testing.T) {
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

func TestBinaryPriorityMap_Pop(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	pm.Set(1, 50)
	pm.Set(2, 10)
	pm.Set(3, 30)

	k, v, ok := pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, 2, k)
	assert.Equal(t, 10, v)

	pm.Pop()
	pm.Pop()

	_, _, ok = pm.Pop()
	assert.False(t, ok)
}

func TestBinaryPriorityMap_Peek(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())

	_, _, ok := pm.Peek()
	assert.False(t, ok)

	pm.Set("top", 1)
	k, v, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "top", k)
	assert.Equal(t, 1, v)
}

func TestBinaryPriorityMap_Contains(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
	pm.Set("a", 10)

	assert.False(t, !pm.Contains("a"))
	assert.True(t, pm.Contains("a"))
	assert.False(t, pm.Contains("b"))

	pm.Remove("a")
	assert.False(t, pm.Contains("a"))
}

func TestBinaryPriorityMap_IsEmpty(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	assert.True(t, pm.IsEmpty())

	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())

	pm.Pop()
	assert.True(t, pm.IsEmpty())
}

func TestBinaryPriorityMap_Length(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, comparator.Min[int]())
	assert.Equal(t, 0, pm.Length())

	pm.Set(1, 10)
	pm.Set(2, 20)
	assert.Equal(t, 2, pm.Length())
}

func TestBinaryPriorityMap_Keys(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.Keys() {
			count++
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 3, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.Keys() {
			count++
			break
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 1, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		count := 0
		for range pm.Keys() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestBinaryPriorityMap_Values(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.Values() {
			count++
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 3, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.Values() {
			count++
			break
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 1, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		count := 0
		for range pm.Values() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestBinaryPriorityMap_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.All() {
			count++
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 3, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("PartialIteration", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		count := 0
		for range pm.All() {
			count++
			break
		}

		assert.Equal(t, 3, pm.Length())
		assert.Equal(t, 1, count)
		assert.True(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		count := 0
		for range pm.Values() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestBinaryPriorityMap_Drain(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
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

	t.Run("PartialIteration", func(t *testing.T) {
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

	t.Run("EmptyMap", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, comparator.Min[int]())
		count := 0
		for range pm.Drain() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestBinaryPriorityMap_Integrity(t *testing.T) {
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

func TestBinaryPriorityMap_Clear(t *testing.T) {
	t.Run("PopulatedMap", func(t *testing.T) {
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

	t.Run("Reuse", func(t *testing.T) {
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
