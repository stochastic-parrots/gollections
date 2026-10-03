package heap

import (
	"cmp"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/heap"
	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

// BinaryHeap is a comparator-backed binary Heap. Its zero value is invalid;
// use NewBinary or NewOrderedBinary.
type BinaryHeap[T any] = heap.BinaryHeap[T]

var _ Heap[int] = &heap.BinaryHeap[int]{}

func checkPriority[T any](hasPriority func(T, T) bool) {
	if hasPriority == nil {
		panic("heap: nil priority comparator")
	}
}

func orderedPriority[T cmp.Ordered](order Order) func(T, T) bool {
	if order == Max {
		return ordering.Max[T]()
	}
	return ordering.Min[T]()
}

// NewBinary creates an empty binary heap with a stable strict-weak-order priority
// predicate. Choose it for O(1) Peek and O(log N) Pop or Replace.
// The predicate must remain stable for the heap's lifetime.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(capacity).
func NewBinary[T any](hasPriority func(T, T) bool, capacity int) *BinaryHeap[T] {
	checkPriority(hasPriority)
	return heap.NewBinaryHeap(capacity, hasPriority)
}

// BinaryFrom heapifies data in place and transfers its backing array.
// The caller must stop using data and every alias afterward. Use BinaryClone to preserve it.
// The predicate must be non-nil and define a stable strict weak ordering.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data)).
func BinaryFrom[T any](hasPriority func(T, T) bool, data []T) *BinaryHeap[T] {
	checkPriority(hasPriority)
	return heap.NewBinaryHeapFromSlice(data, hasPriority)
}

// BinaryClone heapifies an independent shallow copy of data.
// It does not modify or retain data's backing array. The predicate must be non-nil
// and define a stable strict weak ordering.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data)).
func BinaryClone[T any](hasPriority func(T, T) bool, data []T) *BinaryHeap[T] {
	checkPriority(hasPriority)
	return heap.NewBinaryHeapCloneSlice(data, hasPriority)
}

// BinaryFromSeq consumes seq once into new storage and heapifies the collected values.
// The predicate must be non-nil and define a stable strict weak ordering.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func BinaryFromSeq[T any](hasPriority func(T, T) bool, seq iter.Seq[T]) *BinaryHeap[T] {
	checkPriority(hasPriority)
	return heap.NewBinaryHeapFromSlice(slices.Collect(seq), hasPriority)
}

// NewOrderedBinary creates an empty heap using T's natural order in the selected direction.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(capacity).
func NewOrderedBinary[T cmp.Ordered](order Order, capacity int) *BinaryHeap[T] {
	return NewBinary(orderedPriority[T](order), capacity)
}

// OrderedBinaryFrom heapifies data in place using T's natural order and transfers its backing array.
// The caller must stop using data and every alias afterward. Use OrderedBinaryClone to preserve it.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data)).
func OrderedBinaryFrom[T cmp.Ordered](order Order, data []T) *BinaryHeap[T] {
	return BinaryFrom(orderedPriority[T](order), data)
}

// OrderedBinaryClone heapifies an independent shallow copy of data using T's natural order.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(len(data)).
func OrderedBinaryClone[T cmp.Ordered](order Order, data []T) *BinaryHeap[T] {
	return BinaryClone(orderedPriority[T](order), data)
}

// OrderedBinaryFromSeq consumes seq once into new storage and heapifies the collected values.
//
// Push uses the O(N + len(xs)) heapify path when the heap is empty or when
// len(xs) exceeds both N and 64. It uses
// O(len(xs) log (N + len(xs))) insertion otherwise.
//
// Performance Summary (Time Complexity):
//
//	Operation         Time Complexity
//	---------------   ---------------
//	IsEmpty()         O(1)
//	Len()             O(1)
//	All()             O(N)
//	Enumerate()       O(N)
//	Push(xs... T)     O(len(xs) log (N + len(xs))) or O(N + len(xs))
//	Pop()             O(log N)
//	Peek()            O(1)
//	Drain()           O(N log N)
//	Replace(x)        O(log N)
//	Clear()           O(N)
//	MarshalJSON()     O(N)
//	String()          O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func OrderedBinaryFromSeq[T cmp.Ordered](order Order, seq iter.Seq[T]) *BinaryHeap[T] {
	return BinaryFromSeq(orderedPriority[T](order), seq)
}
