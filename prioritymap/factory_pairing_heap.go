package prioritymap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/comparator"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

// PairingHeapPriorityMap is an indexed, pointer-based pairing heap. Its zero
// value is invalid; construct one with [PairingHeap] or [OrderedPairingHeap].
type PairingHeapPriorityMap[K comparable, P any] = prioritymap.PairingPriorityMap[K, P]

var _ PriorityMap[int, any] = &prioritymap.PairingPriorityMap[int, any]{}

// PairingHeapFactory constructs pairing-heap priority maps with a fixed comparator.
//
// Pairing priority maps provide O(1) lookup and Peek. Priority improvements are
// O(1) amortized, while Pop, removal, and arbitrary updates are O(log N)
// amortized.
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
func PairingHeap[K comparable, P any](hasPriority func(P, P) bool) PairingHeapFactory[K, P] {
	return PairingHeapFactory[K, P]{hasPriority: hasPriority}
}

// OrderedPairingHeap returns a pairing-heap factory using the natural order of P.
func OrderedPairingHeap[K comparable, P cmp.Ordered](order Order) PairingHeapFactory[K, P] {
	if order == Max {
		return PairingHeap[K](comparator.Max[P]())
	}
	return PairingHeap[K](comparator.Min[P]())
}

// New creates an empty pairing priority map and retains up to capacity nodes for
// reuse. Once that limit is reached, additional removed nodes are released
// rather than retained in the freelist.
func (factory PairingHeapFactory[K, P]) New(capacity int) *PairingHeapPriorityMap[K, P] {
	return prioritymap.NewPairingPriorityMapWithCapacity[K](capacity, factory.hasPriority)
}
