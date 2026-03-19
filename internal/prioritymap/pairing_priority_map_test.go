package prioritymap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stretchr/testify/assert"
)

func TestNewPairingPriorityMap(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())

	assert.Equal(t, 0, pm.Length())
	assert.True(t, pm.IsEmpty())
	assert.Nil(t, pm.root)
	assert.Nil(t, pm.freelist)
	assert.Empty(t, pm.indexes)
}

func TestNewPairingPriorityMapWithCapacity(t *testing.T) {
	capacity := 10
	comp := comparator.Min[int]()
	pm := NewPairingPriorityMapWithCapacity[string](capacity, comp)

	assert.NotNil(t, pm)
	assert.Nil(t, pm.root)
	assert.NotNil(t, pm.indexes)
	assert.Equal(t, 0, pm.Length())

	count := 0
	current := pm.freelist
	for current != nil {
		count++
		current = current.next
	}
	assert.Equal(t, capacity, count)
	pm.Set("A", 1)

	count = 0
	current = pm.freelist
	for current != nil {
		count++
		current = current.next
	}
	assert.Equal(t, capacity-1, count)
}

func TestPairingPriorityMap_Merge(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	n := &node[string, int]{key: "A", priority: 10}
	assert.Equal(t, n, pm.merge(nil, n))
	assert.Equal(t, n, pm.merge(n, nil))
}

func TestPairingPriorityMap_Cut(t *testing.T) {
	t.Run("MiddleSibling", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 1)
		pm.Set("B", 10)
		pm.Set("C", 11)
		pm.Set("D", 12)

		nodeC := pm.indexes["C"]
		pm.cut(nodeC)

		assert.Nil(t, nodeC.previous)
		assert.Nil(t, nodeC.next)
		assert.Equal(t, 4, pm.Length())
	})
}

func TestPairingPriorityMap_Combine(t *testing.T) {
	t.Run("EvenChildren", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		// Insert in order
		pm.Set(1, 10)
		pm.Set(2, 20)
		pm.Set(3, 30)
		pm.Set(4, 40)
		pm.Set(5, 50)

		pm.Pop() // Dispach combine with 4 children (even)
		assert.Equal(t, 4, pm.Length())
	})

	t.Run("OddChildren", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		// Insert in order
		pm.Set(1, 10)
		pm.Set(2, 20)
		pm.Set(3, 30)
		pm.Set(4, 40)

		pm.Pop() // Dispach combine with 3 children (odd)
		assert.Equal(t, 3, pm.Length())
	})

	t.Run("NilOrSingle", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		assert.Nil(t, pm.combine(nil))

		n := &node[int, int]{key: 1, priority: 1}
		assert.Equal(t, n, pm.combine(n))
	})
}

func TestPairingPriorityMap_Set(t *testing.T) {
	t.Run("Insert", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, 5, val)
	})

	t.Run("Root", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2)

		assert.Equal(t, 1, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("Internal", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 11)
		pm.Set("B", 2)

		assert.Equal(t, 2, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("BetterCase", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2)

		assert.Equal(t, 1, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("WorseCase", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 2)
		pm.Set("B", 5)
		pm.Set("A", 10)

		_, p, _ := pm.Peek()
		assert.Equal(t, 5, p)
	})

	t.Run("WorseCaseNonRoot", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("Root", 1)
		pm.Set("B", 50)
		pm.Set("A", 10)
		pm.Set("B", 60)

		nodeA := pm.indexes["A"]

		if nodeA.child == nil {
			nodeB := pm.indexes["B"]
			pm.cut(nodeB)
			nodeA.child = nodeB
			nodeB.previous = nodeA
		}

		assert.NotNil(t, nodeA.child)
		pm.Set("A", 100)
		assert.Equal(t, 3, pm.Length())
		valA, _ := pm.Get("A")
		assert.Equal(t, 100, valA)
	})
}

func TestPairingPriorityMap_Update(t *testing.T) {
	t.Run("Nonexistent", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		assert.False(t, pm.Update("nonexistent", 1))
	})

	t.Run("Existent", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
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

func TestPairingPriorityMap_Get(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestPairingPriorityMap_Remove(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("A", 10)
	pm.Set("B", 20)
	pm.Set("C", 30)

	assert.True(t, pm.Remove("B"))
	assert.Equal(t, 2, pm.Length())
	assert.True(t, pm.Remove("A"))
	assert.Equal(t, 1, pm.Length())
	assert.False(t, pm.Remove("non-existent"))
	assert.Equal(t, 1, pm.Length())
}

func TestPairingPriorityMap_Pop(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
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

func TestPairingPriorityMap_Peek(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())

	_, _, ok := pm.Peek()
	assert.False(t, ok)

	pm.Set("top", 1)
	k, p, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "top", k)
	assert.Equal(t, 1, p)
}

func TestPairingPriorityMap_Contains(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("a", 10)

	assert.False(t, !pm.Contains("a"))
	assert.True(t, pm.Contains("a"))
	assert.False(t, pm.Contains("b"))

	pm.Remove("a")
	assert.False(t, pm.Contains("a"))
}

func TestPairingPriorityMap_IsEmpty(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
	assert.True(t, pm.IsEmpty())
	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())
}

func TestPairingPriorityMap_Length(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
	pm.Set(1, 10)
	assert.Equal(t, 1, pm.Length())
}

func TestPairingPriorityMap_Keys(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		count := 0
		for range pm.Keys() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestPairingPriorityMap_Values(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		count := 0
		for range pm.Values() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestPairingPriorityMap_All(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		count := 0
		for range pm.Values() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestPairingPriorityMap_Drain(t *testing.T) {
	t.Run("FullIteration", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Max[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
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
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		count := 0
		for range pm.Drain() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestPairingPriorityMap_Integrity(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
	for i := 10; i > 0; i-- {
		pm.Set(i, i)
	}

	for key, node := range pm.indexes {
		assert.Equal(t, key, node.key)
		if node != pm.root {
			assert.NotNil(t, node.previous)
		}
	}

	for value := range pm.Length() {
		key, priority, ok := pm.Pop()

		assert.Equal(t, value+1, key)
		assert.Equal(t, value+1, priority)
		assert.True(t, ok)
	}
}

func TestPairingPriorityMap_Clear(t *testing.T) {
	t.Run("PopulatedMap", func(t *testing.T) {
		pm := NewPairingPriorityMapWithCapacity[string](10, comparator.Min[int]())
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
		assert.Nil(t, pm.root)

		if pm.freelist != nil {
			curr := pm.freelist
			for curr != nil {
				assert.Nil(t, curr.child)
				assert.Nil(t, curr.next)
				assert.Nil(t, curr.previous)
				curr = curr.next
			}
		}
	})

	t.Run("Reuse", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())

		pm.Set("old", 100)
		pm.Clear()

		pm.Set("new", 5)
		pm.Set("newer", 1)

		assert.Equal(t, 2, pm.Length())
		key, priority, ok := pm.Peek()
		assert.True(t, ok)
		assert.Equal(t, "newer", key)
		assert.Equal(t, 1, priority)

		if pm.freelist != nil {
			curr := pm.freelist
			for curr != nil {
				assert.Nil(t, curr.child)
				assert.Nil(t, curr.next)
				assert.Nil(t, curr.previous)
				curr = curr.next
			}
		}
	})
}
