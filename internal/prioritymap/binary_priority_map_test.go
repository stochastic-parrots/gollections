package prioritymap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBinaryPriorityMap(t *testing.T) {
	pm := NewBinaryPriorityMap[string](10, MinFunc[int]())

	assert.Equal(t, 0, pm.Length())
	assert.True(t, pm.IsEmpty())
	assert.Empty(t, pm.data)
	assert.Empty(t, pm.indexes)
}

func TestBinaryPriorityMapSet(t *testing.T) {
	t.Run("Insert New Keys", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
		pm.Set("A", 10)
		pm.Set("B", 5)

		assert.Equal(t, 2, pm.Length())
		key, val, _ := pm.Peek()
		assert.Equal(t, "B", key)
		assert.Equal(t, 5, val)
	})

	t.Run("Update Existing Keys", func(t *testing.T) {
		pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
		pm.Set("A", 10)
		pm.Set("A", 2) // Update priority to be higher

		assert.Equal(t, 1, pm.Length())
		_, priority, _ := pm.Peek()
		assert.Equal(t, 2, priority)
	})
}

func TestBinaryPriorityMapGet(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
	pm.Set("apple", 100)

	val, ok := pm.Get("apple")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	val, ok = pm.Get("non-existent")
	assert.False(t, ok)
	assert.Zero(t, val)
}

func TestBinaryPriorityMapRemove(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
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
	pm := NewBinaryPriorityMap[int](0, MinFunc[int]())
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
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())

	_, _, ok := pm.Peek()
	assert.False(t, ok)

	pm.Set("top", 1)
	k, v, ok := pm.Peek()
	assert.True(t, ok)
	assert.Equal(t, "top", k)
	assert.Equal(t, 1, v)
}

func TestBinaryPriorityMapIsEmpty(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, MinFunc[int]())
	assert.True(t, pm.IsEmpty())

	pm.Set(1, 10)
	assert.False(t, pm.IsEmpty())

	pm.Pop()
	assert.True(t, pm.IsEmpty())
}

func TestBinaryPriorityMapLength(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, MinFunc[int]())
	assert.Equal(t, 0, pm.Length())

	pm.Set(1, 10)
	pm.Set(2, 20)
	assert.Equal(t, 2, pm.Length())
}

func TestBinaryPriorityMapKeys(t *testing.T) {
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
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
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
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
	pm := NewBinaryPriorityMap[string](0, MinFunc[int]())
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

func TestBinaryPriorityMapIntegrity(t *testing.T) {
	pm := NewBinaryPriorityMap[int](0, MinFunc[int]())

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
