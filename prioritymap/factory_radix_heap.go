package prioritymap

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

// RadixHeapPriorityMap is an indexed monotone min-priority map. Its zero value
// is invalid; construct one with [RadixHeap].
type RadixHeapPriorityMap[K comparable, P constraint.Integer] = prioritymap.RadixPriorityMap[K, P]

var _ PriorityMap[int, uint64] = &prioritymap.RadixPriorityMap[int, uint64]{}

// RadixHeapFactory constructs monotone radix-heap priority maps.
//
// Priorities must be non-negative. After Pop returns priority p, subsequent
// mutations must not introduce a priority smaller than p. This invariant is not
// checked. Set, Update, Improve, Get, and Remove are O(1); Pop is O(W)
// amortized, where W is the priority bit width.
//
// Performance Summary (Amortized Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	New(capacity)       O(capacity + W)
//	Set/Update/Improve  O(1)
//	Get/Remove          O(1)
//	Peek                O(N)
//	Pop                 O(W)
//	Drain               O(NW)
//	Clear               O(N + W)
type RadixHeapFactory[K comparable, P constraint.Integer] struct{}

// RadixHeap returns a factory for monotone radix-heap priority maps.
func RadixHeap[K comparable, P constraint.Integer]() RadixHeapFactory[K, P] {
	return RadixHeapFactory[K, P]{}
}

// New creates an empty radix priority map and retains up to capacity entries for
// reuse. Once that limit is reached, additional removed entries are released
// rather than retained in the freelist.
func (RadixHeapFactory[K, P]) New(capacity int) *RadixHeapPriorityMap[K, P] {
	return prioritymap.NewRadixPriorityMap[K, P](capacity)
}
