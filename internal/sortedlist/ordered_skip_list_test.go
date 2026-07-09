package sortedlist

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderedSkipList(t *testing.T) {
	list := NewOrderedSkipList[int]()

	assert.True(t, list.IsEmpty())
	assert.Equal(t, 1, list.level)
	assert.Equal(t, skipListMaxLevel, list.nodes[skipListHead].height)
	assert.Equal(t, skipListNil, list.tail)
	assert.Equal(t, skipListInitialBlockSize, list.nextBlockSize)
	assert.Zero(t, list.freeLength)
	assertOrderedSkipListInvariants(t, list)
}

func TestNewOrderedSkipListFromSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewOrderedSkipListFromSlice(data)

	assert.Equal(t, []int{1, 2, 3}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assert.Len(t, list.nodes, skipListHead+1+len(data))
	assert.Zero(t, list.freeLength)
	assertOrderedSkipListInvariants(t, list)
}

func TestNewOrderedSkipListCloneSlice(t *testing.T) {
	data := []int{3, 1, 2}

	list := NewOrderedSkipListCloneSlice(data)

	assert.Equal(t, []int{3, 1, 2}, data)
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assertOrderedSkipListInvariants(t, list)
}

func TestNewOrderedSkipListFromSeq(t *testing.T) {
	list := NewOrderedSkipListFromSeq(slices.Values([]int{3, 1, 2}))

	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())
	assertOrderedSkipListInvariants(t, list)
}

func TestOrderedSkipList_Behavior(t *testing.T) {
	list := NewOrderedSkipList[int]()
	data := []int{3, 1, 2, 2}

	assert.True(t, list.IsEmpty())

	list.Add(data...)

	assert.Equal(t, []int{3, 1, 2, 2}, data)
	assert.Equal(t, 4, list.Length())
	assert.Equal(t, []int{1, 2, 2, 3}, list.ToSlice())
	assert.Len(t, list.nodes, skipListHead+1+4)
	assert.Zero(t, list.freeLength)
	assert.Equal(t, []int{3, 2, 2, 1}, slices.Collect(list.Backward()))
	assert.Equal(t, 1, list.LowerBound(2))
	assert.Equal(t, 3, list.UpperBound(2))
	start, end := list.EqualRange(2)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	assert.Equal(t, 2, list.Count(2))
	assert.True(t, list.Contains(3))
	assert.False(t, list.Contains(4))

	x, err := list.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, 1, x)

	err = list.Replace(2, 2)
	assert.NoError(t, err)
	assert.ErrorIs(t, list.Replace(2, 5), ErrOrderViolation)

	assert.True(t, list.Remove(2))
	assert.Equal(t, []int{1, 2, 3}, list.ToSlice())

	list.Clear()
	assert.True(t, list.IsEmpty())
	assert.Nil(t, list.ToSlice())
	assertOrderedSkipListInvariants(t, list)
}

func TestOrderedSkipList_Navigation(t *testing.T) {
	list := NewOrderedSkipList[int]()
	list.Add(1, 3, 3, 5)

	idx, ok := list.Find(3)
	assert.True(t, ok)
	assert.Equal(t, 1, idx)

	idx, ok = list.Find(2)
	assert.False(t, ok)
	assert.Equal(t, -1, idx)

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
	assertOrderedSkipListInvariants(t, list)
}

func TestOrderedSkipList_FirstLast(t *testing.T) {
	empty := NewOrderedSkipList[int]()

	first, firstOK := empty.First()
	last, lastOK := empty.Last()

	assert.False(t, firstOK)
	assert.Zero(t, first)
	assert.False(t, lastOK)
	assert.Zero(t, last)

	list := NewOrderedSkipList[int]()
	list.Add(3, 1, 2)

	first, firstOK = list.First()
	last, lastOK = list.Last()

	assert.True(t, firstOK)
	assert.Equal(t, 1, first)
	assert.True(t, lastOK)
	assert.Equal(t, 3, last)
}

func TestOrderedSkipList_Errors(t *testing.T) {
	ordered := NewOrderedSkipList[int]()

	_, err := ordered.Get(0)
	assert.Error(t, err)
	assert.Error(t, ordered.Replace(0, 1))
	assert.False(t, ordered.Remove(1))

	ordered.Add(1, 3, 5)

	assert.ErrorIs(t, ordered.Replace(-1, 1), list.ErrIndexOutOfBound)
	assert.ErrorIs(t, ordered.Replace(3, 1), list.ErrIndexOutOfBound)
	assert.ErrorIs(t, ordered.Replace(1, 0), ErrOrderViolation)
	assert.ErrorIs(t, ordered.Replace(1, 6), ErrOrderViolation)
	assert.Equal(t, []int{1, 3, 5}, ordered.ToSlice())
	assertOrderedSkipListInvariants(t, ordered)
}

func TestOrderedSkipList_RangeAndIterators(t *testing.T) {
	list := NewOrderedSkipList[int]()
	list.Add(1, 3, 3, 5)

	assert.Equal(t, []int{3, 3}, slices.Collect(list.Range(2, 5)))
	assert.Equal(t, []int{3, 3}, slices.Collect(list.Range(3, 4)))
	assert.Empty(t, slices.Collect(list.Range(5, 2)))
	assert.Equal(t, []int{1, 3, 3, 5}, slices.Collect(list.All()))

	indexes := make([]int, 0, list.Length())
	values := make([]int, 0, list.Length())
	for idx, value := range list.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}

	assert.Equal(t, []int{0, 1, 2, 3}, indexes)
	assert.Equal(t, []int{1, 3, 3, 5}, values)

	var stopped []int
	list.Range(2, 5)(func(x int) bool {
		stopped = append(stopped, x)
		return false
	})
	assert.Equal(t, []int{3}, stopped)
}

func TestOrderedSkipList_FreeListReuse(t *testing.T) {
	t.Run("RemoveClearsNode", func(t *testing.T) {
		list := NewOrderedSkipList[string]()
		list.Add("a", "b")
		removed := list.next(skipListHead, 0)

		assert.True(t, list.Remove("a"))
		assert.Zero(t, list.nodes[removed].height)
		assert.Equal(t, removed, list.free)
		assert.Equal(t, 1, list.freeLength)
		assert.Empty(t, list.nodes[removed].value)
		for _, link := range list.links[list.nodes[removed].linkStart : list.nodes[removed].linkStart+list.nodes[removed].linkCapacity] {
			assert.Zero(t, link)
		}
		assertOrderedSkipListInvariants(t, list)
	})

	t.Run("AddReusesNode", func(t *testing.T) {
		list := NewOrderedSkipList[int]()
		list.reserveNodes(1)
		reusable := list.free

		list.Add(1)

		assert.Equal(t, reusable, list.next(skipListHead, 0))
		assert.Equal(t, skipListNil, list.free)
		assert.Zero(t, list.freeLength)
		assert.Equal(t, []int{1}, list.ToSlice())
		assertOrderedSkipListInvariants(t, list)
	})
}

func TestOrderedSkipList_JSON(t *testing.T) {
	list := NewOrderedSkipList[int]()
	list.Add(3, 1, 2)

	data, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(data))

	err = json.Unmarshal([]byte(`[9,7,8]`), list)
	assert.NoError(t, err)
	assert.Equal(t, []int{7, 8, 9}, list.ToSlice())
	assertOrderedSkipListInvariants(t, list)
}

func TestOrderedSkipList_Format(t *testing.T) {
	list := NewOrderedSkipList[int]()
	list.Add(5, 4, 3, 2, 1, 0)

	assert.Equal(t, "[0 1 2 3 4 ...(+1 more)]", list.String())
	assert.Contains(t, fmt.Sprintf("%#v", list), "OrderedSkipList")
}

func assertOrderedSkipListInvariants[T cmp.Ordered](t *testing.T, list *OrderedSkipList[T]) {
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
