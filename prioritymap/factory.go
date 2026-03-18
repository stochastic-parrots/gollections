package prioritymap

import (
	"cmp"

	comparator "github.com/stochastic-parrots/gollections/internal/comparator"
	constructor "github.com/stochastic-parrots/gollections/internal/prioritymap"
)

var _ PriorityMap[int, any] = &constructor.BinaryPriorityMap[int, any]{}
var _ PriorityMap[int, any] = &constructor.PairingPriorityMap[int, any]{}

// BinaryPriorityMap implements a priority-map (priority queue with lookup)
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
type BinaryPriorityMap[K comparable, P any] = *constructor.BinaryPriorityMap[K, P]

// PairingPriorityMap implements a priority-map (priority queue with lookup)
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
type PairingPriorityMap[K comparable, P any] = *constructor.PairingPriorityMap[K, P]

// NewBinary creates and returns a new empty Priority Map (Indexed Binary Heap) with a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(log N)
//	Get(K)              O(1)
//	Remove(K)           O(log N)
//	Pop()               O(log N)
//	Peek()              O(1)
//	Length()            O(1)
func NewBinary[K comparable, V comparable](capacity int, hasPriority func(V, V) bool) BinaryPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, hasPriority)
}

// MinBinary creates and returns a new empty Priority Map (Indexed Min Binary Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(log N)
//	Get(K)              O(1)
//	Remove(K)           O(log N)
//	Pop()               O(log N)
//	Peek()              O(1)
func MinBinary[K comparable, V cmp.Ordered](capacity int) BinaryPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, comparator.Min[V]())
}

// MaxBinary creates and returns a new empty Priority Map (Indexed Max Binary Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(log N)
//	Get(K)              O(1)
//	Remove(K)           O(log N)
//	Pop()               O(log N)
//	Peek()              O(1)
func MaxBinary[K comparable, V cmp.Ordered](capacity int) BinaryPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, comparator.Max[V]())
}

// NewPairing creates and returns a new empty Priority Map (Indexed Pairing Heap) with a custom comparator.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
//	Length()            O(1)
func NewPairing[K comparable, V comparable](capacity int, hasPriority func(V, V) bool) PairingPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, hasPriority)
}

// MinPairing creates and returns a new empty Priority Map (Indexed Min Pairing Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
func MinPairing[K comparable, V cmp.Ordered](capacity int) PairingPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, comparator.Min[V]())
}

// MaxPairing creates and returns a new empty Priority Map (Indexed Max Pairing Heap) for ordered types.
//
// Performance Summary (Time Complexity):
//
//	Operation           Time Complexity
//	-----------------   ---------------
//	Set(K, V)           O(1) Amortized* (O(log N) if priority worsens)
//	Get(K)              O(1)
//	Remove(K)           O(log N) Amortized
//	Pop()               O(log N) Amortized
//	Peek()              O(1)
func MaxPairing[K comparable, V cmp.Ordered](capacity int) PairingPriorityMap[K, V] {
	return constructor.NewPairingPriorityMapWithCapacity[K](capacity, comparator.Max[V]())
}
