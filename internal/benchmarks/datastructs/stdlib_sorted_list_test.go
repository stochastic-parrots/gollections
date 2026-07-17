package datastructs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdSortedList_Add(t *testing.T) {
	list := NewStdSortedList(0)

	list.Add(3)
	list.Add(1)
	list.Add(2)

	assert.Equal(t, []int{1, 2, 3}, list.data)
}

func TestStdSortedList_Adds(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewStdSortedList(0)

		list.Adds()

		assert.Empty(t, list.data)
	})

	t.Run("Single", func(t *testing.T) {
		list := NewStdSortedList(0)

		list.Adds(2)

		assert.Equal(t, []int{2}, list.data)
	})

	t.Run("Many", func(t *testing.T) {
		list := NewStdSortedList(0)
		list.Add(3)

		list.Adds(4, 1, 2)

		assert.Equal(t, []int{1, 2, 3, 4}, list.data)
	})
}
