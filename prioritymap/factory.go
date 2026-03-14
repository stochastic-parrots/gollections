package prioritymap

import (
	"cmp"

	comparator "github.com/stochastic-parrots/gollections/internal/comparator"
	constructor "github.com/stochastic-parrots/gollections/internal/prioritymap"
)

var _ PriorityMap[int, any] = &constructor.BinaryPriorityMap[int, any]{}

// BinaryPriorityMap implements a priority-map (priority queue with lookup)
// built on top of a binary heap and an internal hash map.
//
// This structure allows associating keys with priorities, enabling
// updates or removals of any element in logarithmic time using its key.
type BinaryPriorityMap[K comparable, V any] = *constructor.BinaryPriorityMap[K, V]

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
func NewBinary[K comparable, V comparable](capacity int, cmp func(V, V) bool) BinaryPriorityMap[K, V] {
	return constructor.NewBinaryPriorityMap[K](capacity, cmp)
}

// Min creates and returns a new empty Priority Map (Indexed Min Binary Heap) for ordered types.
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
