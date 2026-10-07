package datastructs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdSortedList_Add(t *testing.T) {
	t.Run("Incremental", func(t *testing.T) {
		list := NewStdSortedList(0)

		list.Add(3)
		list.Add(1)
		list.Add(2)

		assert.Equal(t, []int{1, 2, 3}, list.data)
	})

	t.Run("Batch", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			list := NewStdSortedList(0)

			list.Add()

			assert.Empty(t, list.data)
		})

		t.Run("Single", func(t *testing.T) {
			list := NewStdSortedList(0)

			list.Add(2)

			assert.Equal(t, []int{2}, list.data)
		})

		t.Run("Many", func(t *testing.T) {
			list := NewStdSortedList(0)
			list.Add(3)

			list.Add(4, 1, 2)

			assert.Equal(t, []int{1, 2, 3, 4}, list.data)
		})
	})
}
