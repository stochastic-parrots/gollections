package prioritymap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stretchr/testify/assert"
)

func TestNewRadixPriorityMap(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](10)

	assert.Equal(t, 0, pm.Length())
	assert.True(t, pm.IsEmpty())
	assert.Empty(t, pm.entries)
	assert.Zero(t, pm.last)
	assert.Equal(t, 10, countRadixFree(pm))
	for _, bucket := range pm.buckets {
		assert.Nil(t, bucket)
	}
}

func TestRadixPriorityMap_Set(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, uint64(5), val)
	})

	t.Run("Update", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("A", 10)
		pm.Set("A", 2)

		assert.Equal(t, 1, pm.Length())
		_, priority, _ := pm.Peek()
		assert.Equal(t, uint64(2), priority)
	})

	t.Run("WorsePriority", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("A", 2)
		pm.Set("B", 5)
		pm.Set("A", 10)

		_, priority, _ := pm.Peek()
		assert.Equal(t, uint64(5), priority)
	})
}

func TestRadixPriorityMap_Update(t *testing.T) {
	t.Run("Nonexistent", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		assert.False(t, pm.Update("nonexistent", 1))
	})

	t.Run("Existent", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("apple", 50)
		pm.Set("banana", 30)
		pm.Set("cherry", 10)

		assert.True(t, pm.Update("apple", 5))

		key, priority, ok := pm.Peek()
		assert.Equal(t, "apple", key)
		assert.Equal(t, uint64(5), priority)
		assert.True(t, ok)
	})
}

func TestRadixPriorityMap_Improve(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)

		assert.True(t, pm.Improve("A", 10))
		assert.Equal(t, 1, pm.Length())

		val, ok := pm.Get("A")
		assert.True(t, ok)
		assert.Equal(t, uint64(10), val)
	})

	t.Run("ImprovePriority", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("Root", 1)
		pm.Set("A", 10)
		pm.Set("B", 20)

		assert.True(t, pm.Improve("A", 5))

		key, priority, _ := pm.Peek()
		assert.Equal(t, "Root", key)
		assert.Equal(t, uint64(1), priority)

		assert.True(t, pm.Improve("Root", 0))

		key, priority, _ = pm.Peek()
		assert.Equal(t, "Root", key)
		assert.Equal(t, uint64(0), priority)
		priority, _ = pm.Get("A")
		assert.Equal(t, uint64(5), priority)
	})

	t.Run("WorsenPriority", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("A", 10)

		assert.False(t, pm.Improve("A", 15))

		val, _ := pm.Get("A")
		assert.Equal(t, uint64(10), val)
	})

	t.Run("SamePriority", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("A", 10)

		assert.False(t, pm.Improve("A", 10))
	})
}

func TestRadixPriorityMap_Get(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](0)
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, uint64(100), val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestRadixPriorityMap_Remove(t *testing.T) {
	t.Run("ExistingAndMissing", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](3)
		pm.Set("A", 10)
		pm.Set("B", 20)
		pm.Set("C", 30)

		assert.True(t, pm.Remove("B"))
		assert.Equal(t, 2, pm.Length())
		assert.False(t, pm.Contains("B"))
		assert.Equal(t, 1, countRadixFree(pm))

		assert.True(t, pm.Remove("A"))
		assert.Equal(t, 1, pm.Length())
		assert.False(t, pm.Remove("non-existent"))
		assert.Equal(t, 1, pm.Length())
	})

	t.Run("HeadWithNext", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("tail", 10)
		pm.Set("head", 12)

		assert.True(t, pm.Remove("head"))
		assert.Equal(t, 1, pm.Length())
		assert.True(t, pm.Contains("tail"))
		assert.False(t, pm.Contains("head"))
	})
}

func TestRadixPriorityMap_Pop(t *testing.T) {
	t.Run("PriorityOrder", func(t *testing.T) {
		pm := NewRadixPriorityMap[int, uint64](0)
		pm.Set(1, 50)
		pm.Set(2, 10)
		pm.Set(3, 30)

		k, v, ok := pm.Pop()
		assert.True(t, ok)
		assert.Equal(t, 2, k)
		assert.Equal(t, uint64(10), v)

		pm.Pop()
		pm.Pop()

		_, _, ok = pm.Pop()
		assert.False(t, ok)
	})

	t.Run("BucketZero", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("zero", 0)

		key, priority, ok := pm.Pop()
		assert.True(t, ok)
		assert.Equal(t, "zero", key)
		assert.Equal(t, uint64(0), priority)
	})
}

func TestRadixPriorityMap_Peek(t *testing.T) {
	t.Run("EmptyAndPopulated", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)

		_, _, ok := pm.Peek()
		assert.False(t, ok)

		pm.Set("top", 1)
		k, p, ok := pm.Peek()
		assert.True(t, ok)
		assert.Equal(t, "top", k)
		assert.Equal(t, uint64(1), p)
	})

	t.Run("ScansBucketForMinimum", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("minimum", 10)
		pm.Set("later", 12)

		key, priority, ok := pm.Peek()
		assert.True(t, ok)
		assert.Equal(t, "minimum", key)
		assert.Equal(t, uint64(10), priority)
	})
}

func TestRadixPriorityMap_PeekDoesNotAdvanceLastPriority(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](0)
	pm.Set("later", 20)

	key, priority, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "later", key)
	assert.Equal(t, uint64(20), priority)
	assert.Zero(t, pm.LastPriority())

	pm.Set("earlier", 10)
	_, priority, ok = pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, uint64(10), priority)
}

func TestRadixPriorityMap_Contains(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](0)
	pm.Set("a", 10)

	assert.True(t, pm.Contains("a"))
	assert.False(t, pm.Contains("b"))

	pm.Remove("a")
	assert.False(t, pm.Contains("a"))
}

func TestRadixPriorityMap_IsEmpty(t *testing.T) {
	pm := NewRadixPriorityMap[int, uint64](0)
	assert.True(t, pm.IsEmpty())

	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())
}

func TestRadixPriorityMap_Length(t *testing.T) {
	pm := NewRadixPriorityMap[int, uint64](0)
	pm.Set(1, 10)

	assert.Equal(t, 1, pm.Length())
}

func TestRadixPriorityMap_LastPriority(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](0)
	assert.Zero(t, pm.LastPriority())

	pm.Set("middle", 50)
	pm.Set("near", 10)

	_, priority, ok := pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, uint64(10), priority)
	assert.Equal(t, uint64(10), pm.LastPriority())

	_, priority, ok = pm.Pop()
	assert.True(t, ok)
	assert.Equal(t, uint64(50), priority)
	assert.Equal(t, uint64(50), pm.LastPriority())
}

func TestRadixPriorityMap_Keys(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
		count := 0
		for range pm.Keys() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestRadixPriorityMap_Values(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
		count := 0
		for range pm.Values() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestRadixPriorityMap_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
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
		pm := NewRadixPriorityMap[string, uint64](0)
		count := 0
		for range pm.All() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestRadixPriorityMap_Drain(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		items := map[string]uint64{"a": 30, "b": 10, "c": 20}
		for k, v := range items {
			pm.Set(k, v)
		}

		keys := []string{"b", "c", "a"}
		priorities := []uint64{10, 20, 30}
		idx := 0
		for key, priority := range pm.Drain() {
			assert.Equal(t, keys[idx], key)
			assert.Equal(t, priorities[idx], priority)
			idx++
		}

		assert.True(t, pm.IsEmpty())
	})

	t.Run("PartialIteration", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("a", 10)
		pm.Set("b", 20)
		pm.Set("c", 30)

		for key := range pm.Drain() {
			assert.Equal(t, "a", key)
			break
		}

		assert.Equal(t, 2, pm.Length())
		assert.False(t, pm.Contains("a"))
		assert.True(t, pm.Contains("b"))
		assert.True(t, pm.Contains("c"))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		count := 0
		for range pm.Drain() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestRadixPriorityMap_Integrity(t *testing.T) {
	pm := NewRadixPriorityMap[int, uint64](0)
	for i := 10; i > 0; i-- {
		pm.Set(i, uint64(i))
	}

	for value := range pm.Length() {
		key, priority, ok := pm.Pop()

		assert.Equal(t, value+1, key)
		assert.Equal(t, uint64(value+1), priority)
		assert.True(t, ok)
	}
}

func TestRadixPriorityMap_FreelistReuse(t *testing.T) {
	pm := NewRadixPriorityMap[string, uint64](2)
	assert.Equal(t, 2, countRadixFree(pm))

	pm.Set("a", 1)
	pm.Set("b", 2)
	assert.Equal(t, 0, countRadixFree(pm))

	assert.True(t, pm.Remove("a"))
	assert.Equal(t, 1, countRadixFree(pm))

	pm.Set("c", 3)
	assert.Equal(t, 0, countRadixFree(pm))
}

func TestRadixPriorityMap_Clear(t *testing.T) {
	t.Run("PopulatedMap", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](10)
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

		assert.Equal(t, 0, len(pm.entries))
		assert.Equal(t, 10, countRadixFree(pm))
		for _, bucket := range pm.buckets {
			assert.Nil(t, bucket)
		}
		for curr := pm.free; curr != nil; curr = curr.next {
			assert.Zero(t, curr.key)
			assert.Zero(t, curr.priority)
			assert.Zero(t, curr.bucket)
			assert.Nil(t, curr.previous)
		}
	})

	t.Run("ResetsLastPriority", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)
		pm.Set("high", 100)
		_, _, ok := pm.Pop()
		assert.True(t, ok)
		assert.Equal(t, uint64(100), pm.LastPriority())

		pm.Clear()

		assert.Zero(t, pm.LastPriority())

		pm.Set("low", 1)
		key, priority, ok := pm.Pop()
		assert.True(t, ok)
		assert.Equal(t, "low", key)
		assert.Equal(t, uint64(1), priority)
	})

	t.Run("Reuse", func(t *testing.T) {
		pm := NewRadixPriorityMap[string, uint64](0)

		pm.Set("old", 100)
		pm.Clear()
		assert.Equal(t, 0, len(pm.entries))
		assert.Equal(t, 1, countRadixFree(pm))

		pm.Set("new", 5)
		pm.Set("newer", 1)

		assert.Equal(t, 2, pm.Length())
		key, priority, ok := pm.Peek()
		assert.True(t, ok)
		assert.Equal(t, "newer", key)
		assert.Equal(t, uint64(1), priority)
	})
}

func countRadixFree[K comparable, P constraint.Integer](pm *RadixPriorityMap[K, P]) int {
	count := 0
	for entry := pm.free; entry != nil; entry = entry.next {
		count++
	}
	return count
}
