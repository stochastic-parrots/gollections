package prioritymap

import (
	"cmp"

	comparator "github.com/stochastic-parrots/gollections/internal/comparator"
	constructor "github.com/stochastic-parrots/gollections/internal/prioritymap"
)

var _ PriorityMap[int, any] = &constructor.BinaryPriorityMap[int, any]{}
var _ PriorityMap[int, any] = &constructor.PairingPriorityMap[int, any]{}

// BinaryHeapPriorityMap implements a priority-map (priority queue with lookup)
// built on top of a classic Binary Heap and an internal hash map.
//
// Unlike pointer-based heaps, this implementation stores its elements in a
// contiguous slice (array-based heap). This provides excellent cache locality
// and predictable O(log N) performance for all mutation operations.
//
// It is the ideal choice when memory overhead must be kept to a minimum
// or when the workload involves a balanced mix of priority improvements
// and worsenings, as its worst-case performance is more stable than
// amortized structures.
type BinaryHeapPriorityMap[K comparable, P any] = *constructor.BinaryPriorityMap[K, P]

// PairingHeapPriorityMap implements a priority-map (priority queue with lookup)
// built on top of a Pairing Heap and an internal hash map.
//
// Unlike the Binary Heap which uses a slice-based array, the Pairing Heap
// is a pointer-based multi-way tree structure. It is renowned for its
// empirical performance, often outperforming Binary and Fibonacci Heaps
// in scenarios with frequent priority updates (Decrease-Key).
//
// This implementation uses a "two-pass" merging strategy during Pop operations,
// which maintains a remarkably flat tree structure, ensuring efficient
// future operations.
type PairingHeapPriorityMap[K comparable, P any] = *constructor.PairingPriorityMap[K, P]

// NewBinary creates and returns a new empty Priority Map (Indexed Binary Heap) with a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//  Set(K, V)           O(log N)
//  Update(K, V)        O(log N)
//  Get(K)              O(1)
//  Remove(K)           O(log N)
//  Pop()               O(log N)
//  Peek()              O(1)
//  Clear()             O(N)
//  Drain()             O(N log N)
func NewBinaryHeap[K comparable, V any](capacity int, hasPriority func(V, V) bool) BinaryHeapPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, hasPriority)
}

// MinBinaryHeap creates and returns a new empty Priority Map (Indexed Min Binary Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//  Set(K, V)           O(log N)
//  Update(K, V)        O(log N)
//  Get(K)              O(1)
//  Remove(K)           O(log N)
//  Pop()               O(log N)
//  Peek()              O(1)
//  Clear()             O(N)
//  Drain()             O(N log N)
func MinBinaryHeap[K comparable, V cmp.Ordered](capacity int) BinaryHeapPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, comparator.Min[V]())
}

// MaxBinaryHeap creates and returns a new empty Priority Map (Indexed Max Binary Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//  Set(K, V)           O(log N)
//  Update(K, V)        O(log N)
//  Get(K)              O(1)
//  Remove(K)           O(log N)
//  Pop()               O(log N)
//  Peek()              O(1)
//  Clear()             O(N)
//  Drain()             O(N log N)
func MaxBinaryHeap[K comparable, V cmp.Ordered](capacity int) BinaryHeapPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, comparator.Max[V]())
}

// NewPairing creates and returns a new empty Priority Map (Indexed Pairing Heap) with a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Update(K)           O(1) Amortized
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
// 	Drain()             O(N log N) Amortized
//	Clear()             O(N)
//
// * Note: Update is O(1) for priority improvements (e.g., decreasing key in a Min-Heap).
// If the priority is worsened, it performs as O(log N).
func NewPairingHeap[K comparable, V any](capacity int, hasPriority func(V, V) bool) PairingHeapPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, hasPriority)
}

// MinPairingHeap creates and returns a new empty Priority Map (Indexed Min Pairing Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Update(K)           O(1) Amortized
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
// 	Drain()             O(N log N) Amortized
//	Clear()             O(N)
//
// * Note: Update is O(1) for priority improvements (e.g., decreasing key in a Min-Heap).
// If the priority is worsened, it performs as O(log N).
func MinPairingHeap[K comparable, V cmp.Ordered](capacity int) PairingHeapPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, comparator.Min[V]())
}

// MaxPairingHeap creates and returns a new empty Priority Map (Indexed Max Pairing Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Update(K)           O(1) Amortized
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
// 	Drain()             O(N log N) Amortized
//	Clear()             O(N)
//
// * Note: Update is O(1) for priority improvements (e.g., decreasing key in a Min-Heap).
// If the priority is worsened, it performs as O(log N).
func MaxPairingHeap[K comparable, V cmp.Ordered](capacity int) PairingHeapPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, comparator.Max[V]())
}
