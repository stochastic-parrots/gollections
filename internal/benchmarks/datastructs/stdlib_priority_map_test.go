package datastructs

import (
	"cmp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdPriorityMapClearResetsLength(t *testing.T) {
	pm := NewStdPriorityMap[string](4, cmp.Less[int])
	pm.Set("a", 2)
	pm.Set("b", 1)

	pm.Clear()

	assert.Equal(t, 0, pm.Length())

	_, _, ok := pm.Pop()
	assert.False(t, ok)

	pm.Set("c", 0)
	key, priority, ok := pm.Pop()

	assert.True(t, ok)
	assert.Equal(t, "c", key)
	assert.Equal(t, 0, priority)
}
