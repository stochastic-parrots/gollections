package prioritymap

import (
	"cmp"

	"github.com/stochastic-parrots/gollections/internal/prioritymap"
	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

// PairingHeapPriorityMap is an indexed, pointer-based pairing heap. Its zero value is invalid;
// use NewPairingHeap or NewOrderedPairingHeap.
type PairingHeapPriorityMap[K comparable, P any] = prioritymap.PairingPriorityMap[K, P]

var _ PriorityMap[int, any] = &prioritymap.PairingPriorityMap[int, any]{}

// NewPairingHeap creates an indexed pairing-heap priority map with a stable
// strict-weak-order predicate. Choose it when priority improvements are frequent.
// It panics for a nil predicate. The freelist retains at most capacity nodes;
// zero capacity disables freelist retention.
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
//	Set(key, priority)       O(log N) Amortized
//	Update(key, priority)    O(log N) Amortized
//	Improve(key, priority)   O(1) Amortized
//	Remove(key)              O(log N) Amortized
//	Pop()                    O(log N) Amortized
//	Drain()                  O(N log N) Amortized
//	Clear()                  O(N)
//
// Complexity: O(capacity).
func NewPairingHeap[K comparable, P any](hasPriority func(P, P) bool, capacity int) *PairingHeapPriorityMap[K, P] {
	if hasPriority == nil {
		panic("prioritymap: nil priority comparator")
	}
	return prioritymap.NewPairingPriorityMapWithCapacity[K](capacity, hasPriority)
}

// NewOrderedPairingHeap creates an empty priority map using P's natural order.
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
//	Set(key, priority)       O(log N) Amortized
//	Update(key, priority)    O(log N) Amortized
//	Improve(key, priority)   O(1) Amortized
//	Remove(key)              O(log N) Amortized
//	Pop()                    O(log N) Amortized
//	Drain()                  O(N log N) Amortized
//	Clear()                  O(N)
//
// Complexity: O(capacity).
func NewOrderedPairingHeap[K comparable, P cmp.Ordered](order Order, capacity int) *PairingHeapPriorityMap[K, P] {
	if order == Max {
		return NewPairingHeap[K](ordering.Max[P](), capacity)
	}
	return NewPairingHeap[K](ordering.Min[P](), capacity)
}
