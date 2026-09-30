package prioritymap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/prioritymap"
	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

// BinaryHeapPriorityMap is an indexed, slice-backed binary heap. Its zero value is invalid;
// use NewBinaryHeap or NewOrderedBinaryHeap.
type BinaryHeapPriorityMap[K comparable, P any] = prioritymap.BinaryPriorityMap[K, P]

var _ PriorityMap[int, any] = &prioritymap.BinaryPriorityMap[int, any]{}

// NewBinaryHeap creates an indexed binary-heap priority map with a stable
// strict-weak-order predicate. Choose it for predictable O(log N) updates.
// It panics for a nil predicate. Lookup and Peek are O(1); updates and Pop are O(log N).
//
// Performance Summary (Time Complexity):
//
//	Operation                Time Complexity
//	----------------------   ---------------
//	Get(key)                 O(1)
//	Contains(key)            O(1)
//	Keys()                   O(N)
//	Values()                 O(N)
//	All()                    O(N)
//	IsEmpty()                O(1)
//	Length()                 O(1)
//	Peek()                   O(1)
//	Set(key, priority)       O(log N)
//	Update(key, priority)    O(log N)
//	Improve(key, priority)   O(log N)
//	Remove(key)              O(log N)
//	Pop()                    O(log N)
//	Drain()                  O(N log N)
//	Clear()                  O(N)
//
// Complexity: O(capacity).
func NewBinaryHeap[K comparable, P any](hasPriority func(P, P) bool, capacity int) *BinaryHeapPriorityMap[K, P] {
	if hasPriority == nil {
		panic("prioritymap: nil priority comparator")
	}
	return prioritymap.NewBinaryPriorityMap[K](capacity, hasPriority)
}

// NewOrderedBinaryHeap creates an empty priority map using P's natural order.
//
// Performance Summary (Time Complexity):
//
//	Operation                Time Complexity
//	----------------------   ---------------
//	Get(key)                 O(1)
//	Contains(key)            O(1)
//	Keys()                   O(N)
//	Values()                 O(N)
//	All()                    O(N)
//	IsEmpty()                O(1)
//	Length()                 O(1)
//	Peek()                   O(1)
//	Set(key, priority)       O(log N)
//	Update(key, priority)    O(log N)
//	Improve(key, priority)   O(log N)
//	Remove(key)              O(log N)
//	Pop()                    O(log N)
//	Drain()                  O(N log N)
//	Clear()                  O(N)
//
// Complexity: O(capacity).
func NewOrderedBinaryHeap[K comparable, P cmp.Ordered](order Order, capacity int) *BinaryHeapPriorityMap[K, P] {
	if order == Max {
		return NewBinaryHeap[K](ordering.Max[P](), capacity)
	}
	return NewBinaryHeap[K](ordering.Min[P](), capacity)
}
