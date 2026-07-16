package prioritymap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/prioritymap"
	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

// PairingHeapPriorityMap is an indexed, pointer-based pairing heap. Its zero
// value is invalid; construct one with [PairingHeap] or [OrderedPairingHeap].
type PairingHeapPriorityMap[K comparable, P any] = prioritymap.PairingPriorityMap[K, P]

var _ PriorityMap[int, any] = &prioritymap.PairingPriorityMap[int, any]{}

// PairingHeapFactory constructs pairing-heap priority maps with a fixed ordering.
//
// Pairing priority maps provide O(1) lookup and Peek. Priority improvements are
// O(1) amortized, while Pop, removal, and arbitrary updates are O(log N)
// amortized.
// The comparator must be non-nil, define a strict weak ordering, and remain
// stable for the lifetime of every map created by the factory.
//
// The zero value is invalid. Create a factory with [PairingHeap] or
// [OrderedPairingHeap].
//
// Performance Summary (Amortized Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity)
//	Improve             O(1) when priority improves
//	Update              O(log N)
//	Get/Peek            O(1)
//	Remove/Pop          O(log N)
//	Drain               O(N log N)
//	Clear               O(N)
type PairingHeapFactory[K comparable, P any] struct {
	hasPriority func(P, P) bool
}

// PairingHeap returns a pairing-heap priority-map factory using hasPriority for ordering.
// It panics if hasPriority is nil.
func PairingHeap[K comparable, P any](hasPriority func(P, P) bool) PairingHeapFactory[K, P] {
	if hasPriority == nil {
		panic("prioritymap: nil priority comparator")
	}
	return PairingHeapFactory[K, P]{hasPriority: hasPriority}
}

// OrderedPairingHeap returns a pairing-heap factory using the natural order of P.
func OrderedPairingHeap[K comparable, P cmp.Ordered](order Order) PairingHeapFactory[K, P] {
	if order == Max {
		return PairingHeap[K](ordering.Max[P]())
	}
	return PairingHeap[K](ordering.Min[P]())
}

// New creates an empty pairing priority map with the requested initial capacity.
// The freelist retains at most capacity nodes for reuse. A zero capacity
// disables freelist retention.
func (factory PairingHeapFactory[K, P]) New(capacity int) *PairingHeapPriorityMap[K, P] {
	return prioritymap.NewPairingPriorityMapWithCapacity[K](capacity, factory.hasPriority)
}
