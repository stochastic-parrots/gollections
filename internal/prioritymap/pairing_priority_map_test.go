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

func TestPairingPriorityMapMerge(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	n := &node[string, int]{key: "A", priority: 10}
	assert.Equal(t, n, pm.merge(nil, n))
	assert.Equal(t, n, pm.merge(n, nil))
}

func TestPairingPriorityMapCut(t *testing.T) {
	t.Run("Cut Middle Sibling", func(t *testing.T) {
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

func TestPairingPriorityMapCombine(t *testing.T) {
	t.Run("Combine Even Number of Childrens", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		// Insert in order
		pm.Set(1, 10)
		pm.Set(2, 20)
		pm.Set(3, 30)
		pm.Set(4, 40)
		pm.Set(5, 50)

		pm.Pop() // Dispach combine with 4 childrens (even)
		assert.Equal(t, 4, pm.Length())
	})

	t.Run("Combine Odd Number of Childrens", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		// Insert in order
		pm.Set(1, 10)
		pm.Set(2, 20)
		pm.Set(3, 30)
		pm.Set(4, 40)

		pm.Pop() // Dispach combine with 3 childrens (odd)
		assert.Equal(t, 3, pm.Length())
	})

	t.Run("Combine Nil or Single", func(t *testing.T) {
		pm := NewPairingPriorityMap[int](comparator.Min[int]())
		assert.Nil(t, pm.combine(nil))

		n := &node[int, int]{key: 1, priority: 1}
		assert.Equal(t, n, pm.combine(n))
	})
}

func TestPairingPriorityMapSet(t *testing.T) {
	t.Run("Insert New Keys", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, 5, val)
	})

	t.Run("Update Root Priority", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2)

		assert.Equal(t, 1, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("Update Some Node Priority", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("B", 11)
		pm.Set("B", 2)

		assert.Equal(t, 2, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("Update Priority Better", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 10)
		pm.Set("A", 2)

		assert.Equal(t, 1, pm.Length())
		_, p, _ := pm.Peek()
		assert.Equal(t, 2, p)
	})

	t.Run("Update Priority Worse", func(t *testing.T) {
		pm := NewPairingPriorityMap[string](comparator.Min[int]())
		pm.Set("A", 2)
		pm.Set("B", 5)
		pm.Set("A", 10)

		_, p, _ := pm.Peek()
		assert.Equal(t, 5, p)
	})

	t.Run("Update Priority Worse (Non-Root with Children)", func(t *testing.T) {
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

func TestPairingPriorityMapGet(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestPairingPriorityMapRemove(t *testing.T) {
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

func TestPairingPriorityMapPop(t *testing.T) {
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

func TestPairingPriorityMapPeek(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())

	_, _, ok := pm.Peek()
	assert.False(t, ok)

	pm.Set("top", 1)
	k, p, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "top", k)
	assert.Equal(t, 1, p)
}

func TestPairingPriorityMapIsEmpty(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
	assert.True(t, pm.IsEmpty())
	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())
}

func TestPairingPriorityMapLength(t *testing.T) {
	pm := NewPairingPriorityMap[int](comparator.Min[int]())
	pm.Set(1, 10)
	assert.Equal(t, 1, pm.Length())
}

func TestPairingPriorityMapKeys(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("A", 1)
	pm.Set("B", 0)

	keys := make(map[string]bool)
	for k := range pm.Keys() {
		keys[k] = true
	}
	assert.True(t, keys["A"])
	assert.True(t, keys["B"])
	assert.False(t, keys["C"])
}

func TestPairingPriorityMapValues(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("A", 100)

	var values []int
	for v := range pm.Values() {
		values = append(values, v)
	}
	assert.Contains(t, values, 100)
	assert.NotContains(t, values, 200)
}

func TestPairingPriorityMapAll(t *testing.T) {
	pm := NewPairingPriorityMap[string](comparator.Min[int]())
	pm.Set("A", 1)

	for k, v := range pm.All() {
		assert.Equal(t, "A", k)
		assert.Equal(t, 1, v)
	}
}

func TestPairingPriorityMapIntegrity(t *testing.T) {
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
