package heap_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/heap"
	"github.com/stretchr/testify/assert"
)

var _ fmt.Stringer = (*heap.BinaryHeap[int])(nil)
var _ fmt.Formatter = (*heap.BinaryHeap[int])(nil)
var _ json.Marshaler = (*heap.BinaryHeap[int])(nil)

func TestConstructorsImplementHeap(t *testing.T) {
	var _ *heap.BinaryHeap[int] = heap.NewBinary(cmp.Less[int], 0)
	var _ heap.Heap[int] = heap.NewBinary(cmp.Less[int], 0)
	var _ heap.Heap[int] = heap.NewOrderedBinary[int](heap.Min, 0)
	var _ heap.Heap[int] = heap.NewOrderedBinary[int](heap.Max, 0)
	var _ heap.Heap[int] = heap.BinaryClone(cmp.Less[int], []int{1})
}

func TestBinaryConstructors_NilComparator(t *testing.T) {
	constructors := map[string]func(){
		"New":     func() { heap.NewBinary[int](nil, 0) },
		"From":    func() { heap.BinaryFrom[int](nil, nil) },
		"Clone":   func() { heap.BinaryClone[int](nil, nil) },
		"FromSeq": func() { heap.BinaryFromSeq[int](nil, slices.Values([]int(nil))) },
	}
	for name, construct := range constructors {
		t.Run(name, func(t *testing.T) {
			assert.PanicsWithValue(t, "heap: nil priority comparator", construct)
		})
	}
}

func TestNewBinary(t *testing.T) {
	t.Run("Custom", func(t *testing.T) {
		assertHeapBehavior(t, heap.NewBinary(cmp.Less[int], 0), []int{1, 2, 3})
	})
	t.Run("Min", func(t *testing.T) {
		assertHeapBehavior(t, heap.NewOrderedBinary[int](heap.Min, 0), []int{1, 2, 3})
	})
	t.Run("Max", func(t *testing.T) {
		assertHeapBehavior(t, heap.NewOrderedBinary[int](heap.Max, 0), []int{3, 2, 1})
	})
}

func TestBinaryFrom(t *testing.T) {
	data := []int{3, 1, 2}
	h := heap.BinaryFrom(cmp.Less[int], data)

	assert.Equal(t, []int{1, 2, 3}, slices.Collect(func(yield func(int) bool) {
		for _, value := range h.Drain() {
			if !yield(value) {
				return
			}
		}
	}))
}

func TestBinaryClone(t *testing.T) {
	data := []int{3, 1, 2}
	h := heap.BinaryClone(cmp.Less[int], data)
	assert.Equal(t, []int{3, 1, 2}, data)
	data[0] = 9
	_, ok := h.Replace(0)
	assert.True(t, ok)
	assert.Equal(t, []int{9, 1, 2}, data)

	assert.Equal(t, []int{0, 2, 3}, slices.Collect(func(yield func(int) bool) {
		for _, value := range h.Drain() {
			if !yield(value) {
				return
			}
		}
	}))
}

func TestBinaryFromSeq(t *testing.T) {
	calls := 0
	seq := func(yield func(int) bool) {
		calls++
		for _, value := range []int{3, 1, 2} {
			if !yield(value) {
				return
			}
		}
	}
	h := heap.BinaryFromSeq(cmp.Less[int], seq)
	assert.Equal(t, 1, calls)
	assert.Equal(t, []int{1, 2, 3}, drain(h))
}

func TestOrderedBinaryFromSeq(t *testing.T) {
	min := heap.OrderedBinaryFromSeq(heap.Min, slices.Values([]int{3, 1, 2}))
	max := heap.OrderedBinaryFromSeq(heap.Max, slices.Values([]int{3, 1, 2}))
	assert.Equal(t, []int{1, 2, 3}, drain(min))
	assert.Equal(t, []int{3, 2, 1}, drain(max))
	assert.True(t, heap.OrderedBinaryFromSeq(heap.Min, slices.Values([]int(nil))).IsEmpty())
}

func TestEmptySliceConstructors(t *testing.T) {
	for name, construct := range map[string]func() heap.Heap[int]{
		"BinaryFromNil":    func() heap.Heap[int] { return heap.BinaryFrom(cmp.Less[int], []int(nil)) },
		"BinaryFromEmpty":  func() heap.Heap[int] { return heap.BinaryFrom(cmp.Less[int], make([]int, 0, 4)) },
		"BinaryCloneNil":   func() heap.Heap[int] { return heap.BinaryClone(cmp.Less[int], []int(nil)) },
		"BinaryCloneEmpty": func() heap.Heap[int] { return heap.BinaryClone(cmp.Less[int], make([]int, 0, 4)) },
	} {
		t.Run(name, func(t *testing.T) {
			values := construct()
			assert.True(t, values.IsEmpty())
			values.Push(1)
			top, ok := values.Peek()
			assert.True(t, ok)
			assert.Equal(t, 1, top)
		})
	}
}

func TestBinaryHeap_FormattingAndJSON(t *testing.T) {
	for _, test := range []struct {
		name       string
		input      []int
		wantJSON   string
		wantString string
	}{
		{"Empty", nil, "[]", "[]"},
		{"HeapOrder", []int{1, 3, 2}, "[1,3,2]", "[1 3 2]"},
		{"Truncated", []int{1, 3, 2, 7, 5, 4}, "[1,3,2,7,5,4]", "[1 3 2 7 5 ...(+1 more)]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var values heap.Heap[int] = heap.OrderedBinaryFrom(heap.Min, slices.Clone(test.input))
			data, err := json.Marshal(values)
			assert.NoError(t, err)
			assert.Equal(t, test.wantJSON, string(data))
			assert.Equal(t, test.wantString, values.(fmt.Stringer).String())
			assert.Equal(t, test.wantString, fmt.Sprint(values))
		})
	}
}

func assertHeapBehavior(t *testing.T, h heap.Heap[int], expectedDrain []int) {
	t.Helper()

	_, unmarshals := any(h).(json.Unmarshaler)
	assert.False(t, unmarshals)
	assert.True(t, h.IsEmpty())

	h.Push(3)
	h.Pushes(1, 2)

	top, ok := h.Peek()
	assert.True(t, ok)
	assert.Equal(t, expectedDrain[0], top)
	assert.Equal(t, 3, h.Len())
	assert.Len(t, slices.Collect(h.All()), 3)

	old, ok := h.Replace(4)
	assert.True(t, ok)
	assert.Equal(t, expectedDrain[0], old)
	assert.Equal(t, 3, h.Len())

	data, err := json.Marshal(h)
	assert.NoError(t, err)

	var values []int
	err = json.Unmarshal(data, &values)
	assert.NoError(t, err)
	assert.Len(t, values, 3)

	h.Clear()
	assert.True(t, h.IsEmpty())

	h.Pushes(3, 1, 2)
	assert.Equal(t, expectedDrain, drain(h))
	assert.True(t, h.IsEmpty())
}

func drain(h heap.Heap[int]) []int {
	var values []int
	for _, value := range h.Drain() {
		values = append(values, value)
	}
	return values
}
