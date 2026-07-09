package sortedlist

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stretchr/testify/assert"
)

func TestNewSkipList(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])

	assert.True(t, list.IsEmpty())
	assert.Equal(t, 1, list.level)
	assert.Equal(t, skipListMaxLevel, list.nodes[skipListHead].height)
	assert.Equal(t, skipListNil, list.tail)
	assert.Equal(t, skipListInitialBlockSize, list.nextBlockSize)
	assert.Zero(t, list.freeLength)
	assertSkipListInvariants(t, list)
}

func TestNewSkipListFromSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewSkipListFromSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Len(t, list.nodes, skipListHead+1+len(data))
	assert.Zero(t, list.freeLength)
	assertSkipListInvariants(t, list)
}

func TestNewSkipListCloneSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewSkipListCloneSlice(data, cmp.Compare[int])

	assert.Equal(t, []int{3, 1, 2}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assertSkipListInvariants(t, list)
}

func TestNewSkipListFromSeq(t *testing.T) {
	list := NewSkipListFromSeq(slices.Values([]int{3, 1, 2}), cmp.Compare[int])

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assertSkipListInvariants(t, list)
}

func TestSkipList_Bounds(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.Add(1, 3, 3, 5)

	assert.Equal(t, 0, list.LowerBound(0))
	assert.Equal(t, 1, list.LowerBound(3))
	assert.Equal(t, 3, list.LowerBound(4))
	assert.Equal(t, 4, list.LowerBound(6))

	assert.Equal(t, 0, list.UpperBound(0))
	assert.Equal(t, 3, list.UpperBound(3))
	assert.Equal(t, 3, list.UpperBound(4))
	assert.Equal(t, 4, list.UpperBound(5))
	assert.Equal(t, 4, list.UpperBound(6))

	start, end := list.EqualRange(3)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	assert.Equal(t, 2, list.Count(3))

	start, end = list.EqualRange(2)
	assert.Equal(t, 1, start)
	assert.Equal(t, 1, end)
	assert.Zero(t, list.Count(2))
	assertSkipListInvariants(t, list)
}

func TestSkipList_Get(t *testing.T) {
	t.Run("ValidIndex", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(3, 1, 2)

		for idx, value := range []int{1, 2, 3} {
			x, err := list.Get(idx)
			assert.NoError(t, err)
			assert.Equal(t, value, x)
		}
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewSkipList(cmp.Compare[int])

		for _, idx := range []int{-1, 0, 1} {
			_, err := l.Get(idx)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, list.ErrIndexOutOfBound))
		}
	})
}

func TestSkipList_FindContainsAndNavigation(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.Add(1, 3, 3, 5)

	idx, ok := list.Find(3)
	assert.True(t, ok)
	assert.Equal(t, 1, idx)

	idx, ok = list.Find(2)
	assert.False(t, ok)
	assert.Equal(t, -1, idx)

	assert.True(t, list.Contains(5))
	assert.False(t, list.Contains(6))

	value, idx, ok := list.Ceiling(2)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 1, idx)

	value, idx, ok = list.Ceiling(6)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = list.Floor(4)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = list.Floor(0)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = list.Higher(3)
	assert.True(t, ok)
	assert.Equal(t, 5, value)
	assert.Equal(t, 3, idx)

	value, idx, ok = list.Higher(5)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)

	value, idx, ok = list.Lower(4)
	assert.True(t, ok)
	assert.Equal(t, 3, value)
	assert.Equal(t, 2, idx)

	value, idx, ok = list.Lower(1)
	assert.False(t, ok)
	assert.Zero(t, value)
	assert.Equal(t, -1, idx)
}

func TestSkipList_FirstLast(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])

		first, firstOK := list.First()
		last, lastOK := list.Last()

		assert.False(t, firstOK)
		assert.Zero(t, first)
		assert.False(t, lastOK)
		assert.Zero(t, last)
	})

	t.Run("NonEmpty", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(3, 1, 2)

		first, firstOK := list.First()
		last, lastOK := list.Last()

		assert.True(t, firstOK)
		assert.Equal(t, 1, first)
		assert.True(t, lastOK)
		assert.Equal(t, 3, last)
	})
}

func TestSkipList_Add(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])

		list.Add()

		assert.Nil(t, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("SingleValue", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(2)
		list.Add(1)
		list.Add(3)

		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("MultipleValues", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		data := []int{3, 1, 2, 2}

		list.Add(data...)

		assert.Equal(t, []int{3, 1, 2, 2}, data)
		assert.Equal(t, []int{1, 2, 2, 3}, list.ToSlice())
		assert.Len(t, list.nodes, skipListHead+1+4)
		assert.Zero(t, list.freeLength)
		assertSkipListInvariants(t, list)
	})
}

func TestSkipList_Replace(t *testing.T) {
	t.Run("PreservesOrder", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(1, 3, 5)

		err := list.Replace(1, 4)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 4, 5}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("EquivalentValue", func(t *testing.T) {
		type item struct {
			priority int
			name     string
		}
		compare := func(a, b item) int {
			return cmp.Compare(a.priority, b.priority)
		}
		items := NewSkipList(compare)
		items.Add(item{priority: 1, name: "old"})

		err := items.Replace(0, item{priority: 1, name: "new"})

		assert.NoError(t, err)
		value, getErr := items.Get(0)
		assert.NoError(t, getErr)
		assert.Equal(t, "new", value.name)
	})

	t.Run("InvalidIndex", func(t *testing.T) {
		l := NewSkipList(cmp.Compare[int])
		l.Add(1)

		for _, idx := range []int{-1, 1} {
			err := l.Replace(idx, 1)
			assert.ErrorIs(t, err, list.ErrIndexOutOfBound)
		}
	})

	t.Run("OrderViolation", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(1, 3, 5)

		assert.ErrorIs(t, list.Replace(1, 0), ErrOrderViolation)
		assert.ErrorIs(t, list.Replace(1, 6), ErrOrderViolation)
		assert.Equal(t, []int{1, 3, 5}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})
}

func TestSkipList_Remove(t *testing.T) {
	t.Run("ElementExists", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(3, 1, 2, 2)

		ok := list.Remove(2)

		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("NonExistent", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(1, 2, 3)

		ok := list.Remove(4)

		assert.False(t, ok)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("Empty", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])

		assert.False(t, list.Remove(1))
		assertSkipListInvariants(t, list)
	})

	t.Run("ClearsDiscardedReference", func(t *testing.T) {
		list := NewSkipList(func(a, b *int) int {
			return cmp.Compare(*a, *b)
		})
		x := 1
		y := 2
		list.Add(&x, &y)
		removed := list.next(skipListHead, 0)

		assert.True(t, list.Remove(&x))
		assert.Zero(t, list.nodes[removed].height)
		assert.Equal(t, removed, list.free)
		assert.Equal(t, 1, list.freeLength)
		var zero *int
		assert.Equal(t, zero, list.nodes[removed].value)
		for _, link := range list.links[list.nodes[removed].linkStart : list.nodes[removed].linkStart+list.nodes[removed].linkCapacity] {
			assert.Zero(t, link)
		}
		assertSkipListInvariants(t, list)
	})
}

func TestSkipList_ReusesFreeNode(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.reserveNodes(1)
	reusable := list.free

	list.Add(1)

	assert.Equal(t, reusable, list.next(skipListHead, 0))
	assert.Equal(t, skipListNil, list.free)
	assert.Zero(t, list.freeLength)
	assert.Equal(t, []int{1}, list.ToSlice())
	assertSkipListInvariants(t, list)
}

func TestSkipList_Iterators(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.Add(3, 1, 2)

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(list.All()))
	assert.Equal(t, []int{3, 2, 1}, slices.Collect(list.Backward()))

	var allValues []int
	list.All()(func(x int) bool {
		allValues = append(allValues, x)
		return false
	})
	assert.Equal(t, []int{1}, allValues)

	var backwardValues []int
	list.Backward()(func(x int) bool {
		backwardValues = append(backwardValues, x)
		return false
	})
	assert.Equal(t, []int{3}, backwardValues)

	indexes := make([]int, 0, list.Length())
	values := make([]int, 0, list.Length())
	for idx, value := range list.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}

	assert.Equal(t, []int{0, 1, 2}, indexes)
	assert.Equal(t, []int{1, 2, 3}, values)

	var stopped []int
	list.Enumerate()(func(_ int, x int) bool {
		stopped = append(stopped, x)
		return false
	})
	assert.Equal(t, []int{1}, stopped)

	empty := NewSkipList(cmp.Compare[int])
	assert.Empty(t, slices.Collect(empty.All()))
	assert.Empty(t, slices.Collect(empty.Backward()))
	count := 0
	for range empty.Enumerate() {
		count++
	}
	assert.Zero(t, count)
}

func TestSkipList_Range(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.Add(1, 3, 3, 5)

	assert.Equal(t, []int{3, 3}, slices.Collect(list.Range(2, 5)))
	assert.Equal(t, []int{3, 3}, slices.Collect(list.Range(3, 4)))
	assert.Empty(t, slices.Collect(list.Range(5, 2)))

	var values []int
	list.Range(2, 5)(func(x int) bool {
		values = append(values, x)
		return false
	})
	assert.Equal(t, []int{3}, values)
}

func TestSkipList_ToSlice(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	assert.Nil(t, list.ToSlice())

	list.Add(2, 1)
	slice := list.ToSlice()
	slice[0] = 99

	assert.Equal(t, []int{1, 2}, list.ToSlice())
}

func TestSkipList_Clear(t *testing.T) {
	list := NewSkipList(func(a, b *int) int {
		return cmp.Compare(*a, *b)
	})
	x := 1
	y := 2
	list.Add(&x, &y)
	first := list.next(skipListHead, 0)

	list.Clear()

	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.ToSlice())
	assert.Equal(t, skipListNil, list.tail)
	assert.Equal(t, 1, list.level)
	assert.Equal(t, skipListNil, list.next(skipListHead, 0))
	assert.NotEqual(t, skipListNil, list.free)
	assert.Equal(t, 2, list.freeLength)
	assert.Nil(t, list.nodes[first].value)
	assert.Zero(t, list.nodes[first].height)
	for _, link := range list.links[list.nodes[first].linkStart : list.nodes[first].linkStart+list.nodes[first].linkCapacity] {
		assert.Zero(t, link)
	}
	assertSkipListInvariants(t, list)
}

func TestSkipList_JSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(3, 1, 2)

		data, err := json.Marshal(list)

		assert.NoError(t, err)
		assert.JSONEq(t, `[1,2,3]`, string(data))
	})

	t.Run("MarshalEmpty", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])

		data, err := json.Marshal(list)

		assert.NoError(t, err)
		assert.JSONEq(t, `[]`, string(data))
	})

	t.Run("Unmarshal", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(10)

		err := json.Unmarshal([]byte(`[3,1,2]`), list)

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		list := NewSkipList(cmp.Compare[int])
		list.Add(1)

		err := json.Unmarshal([]byte(`{}`), list)

		assert.Error(t, err)
		assert.Equal(t, []int{1}, list.ToSlice())
		assertSkipListInvariants(t, list)
	})
}

func TestSkipList_Format(t *testing.T) {
	list := NewSkipList(cmp.Compare[int])
	list.Add(5, 4, 3, 2, 1, 0)

	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", list.String())
	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", fmt.Sprintf("%v", list))
	assert.Contains(t, fmt.Sprintf("%+v", list), "len:6")
	assert.Contains(t, fmt.Sprintf("%#v", list), "SkipList")
}

func TestSkipList_DescendingComparator(t *testing.T) {
	list := NewSkipList(func(a, b int) int {
		return cmp.Compare(b, a)
	})

	list.Add(1, 3, 2)

	assert.Equal(t, []int{3, 2, 1}, list.ToSlice())
	first, ok := list.First()
	assert.True(t, ok)
	assert.Equal(t, 3, first)
	last, ok := list.Last()
	assert.True(t, ok)
	assert.Equal(t, 1, last)
	assert.Equal(t, []int{3, 2}, slices.Collect(list.Range(3, 1)))
	assertSkipListInvariants(t, list)
}

func assertSkipListInvariants[T any](t *testing.T, list *SkipList[T]) {
	t.Helper()

	count := 0
	previous := skipListNil
	for node := list.next(skipListHead, 0); node != skipListNil; node = list.next(node, 0) {
		assert.Equal(t, previous, list.nodes[node].backward)
		previous = node
		count++
	}

	assert.Equal(t, list.length, count)
	assert.Equal(t, previous, list.tail)
	if list.length == 0 {
		assert.Equal(t, skipListNil, list.tail)
	}

	ranks := map[int]int{skipListHead: 0}
	rank := 1
	for node := list.next(skipListHead, 0); node != skipListNil; node = list.next(node, 0) {
		ranks[node] = rank
		rank++
	}

	freeCount := 0
	seenFree := map[int]bool{}
	for node := list.free; node != skipListNil; node = list.nodes[node].freeNext {
		if seenFree[node] {
			t.Fatalf("freelist cycle at node %d", node)
		}
		seenFree[node] = true
		_, active := ranks[node]
		assert.False(t, active)
		assert.Zero(t, list.nodes[node].height)
		freeCount++
	}
	assert.Equal(t, list.freeLength, freeCount)

	for level := range list.level {
		current := skipListHead
		for list.next(current, level) != skipListNil {
			link := list.link(current, level)
			assert.Greater(t, link.span, 0)
			next := link.next
			assert.Equal(t, ranks[next]-ranks[current], link.span)
			current = next
			assert.GreaterOrEqual(t, list.nodes[current].height, level+1)
		}
	}
}
