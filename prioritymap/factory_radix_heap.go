package prioritymap

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

// RadixHeapPriorityMap is an indexed monotone min-priority map. Its zero value is invalid;
// use NewRadixHeap. Priorities must be non-negative, and after Pop returns priority p,
// subsequent mutations must not introduce a priority smaller than p.
type RadixHeapPriorityMap[K comparable, P constraint.Integer] = prioritymap.RadixPriorityMap[K, P]

var _ PriorityMap[int, uint64] = &prioritymap.RadixPriorityMap[int, uint64]{}

// NewRadixHeap creates an empty radix priority map. Its freelist retains at most
// capacity entries; zero capacity disables freelist retention.
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
//	Len()                 O(1)
//	Peek()                   O(N)
//	Set(key, priority)       O(1)
//	Update(key, priority)    O(1)
//	Improve(key, priority)   O(1)
//	Remove(key)              O(1)
//	Pop()                    O(W) Amortized
//	Drain()                  O(NW) Amortized
//	Clear()                  O(N + W)
//
// Complexity: O(capacity + W), where W is the priority bit width.
func NewRadixHeap[K comparable, P constraint.Integer](capacity int) *RadixHeapPriorityMap[K, P] {
	return prioritymap.NewRadixPriorityMap[K, P](capacity)
}
