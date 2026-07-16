package prioritymap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

// BinaryHeapPriorityMap is an indexed, slice-backed binary heap. Its zero value
// is invalid; construct one with [BinaryHeap] or [OrderedBinaryHeap].
type BinaryHeapPriorityMap[K comparable, P any] = prioritymap.BinaryPriorityMap[K, P]

var _ PriorityMap[int, any] = &prioritymap.BinaryPriorityMap[int, any]{}

// BinaryHeapFactory constructs binary-heap priority maps with a fixed comparator.
//
// Binary priority maps provide O(1) lookup and Peek with predictable O(log N)
// insertion, update, removal, and Pop.
//
// The zero value is invalid. Create a factory with [BinaryHeap] or
// [OrderedBinaryHeap].
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	Set/Update/Improve  O(log N)
//	Get/Peek            O(1)
//	Remove/Pop          O(log N)
//	Drain               O(N log N)
//	Clear               O(N)
type BinaryHeapFactory[K comparable, P any] struct {
	hasPriority func(P, P) bool
}

// BinaryHeap returns a binary-heap priority-map factory using hasPriority for ordering.
func BinaryHeap[K comparable, P any](hasPriority func(P, P) bool) BinaryHeapFactory[K, P] {
	return BinaryHeapFactory[K, P]{hasPriority: hasPriority}
}

// OrderedBinaryHeap returns a binary-heap priority-map factory using the natural order of P.
func OrderedBinaryHeap[K comparable, P cmp.Ordered](order Order) BinaryHeapFactory[K, P] {
	if order == Max {
		return BinaryHeap[K](comparator.Max[P]())
	}
	return BinaryHeap[K](comparator.Min[P]())
}

// New creates an empty binary priority map with the requested initial capacity.
func (factory BinaryHeapFactory[K, P]) New(capacity int) *BinaryHeapPriorityMap[K, P] {
	return prioritymap.NewBinaryPriorityMap[K](capacity, factory.hasPriority)
}
