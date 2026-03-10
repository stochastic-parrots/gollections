package heap

import (
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
)

// Heap defines a priority queue structure where elements are ordered
// based on a specific priority criteria (Min-Heap or Max-Heap).
// It extends the basic functionalities provided by the gollections.Collection interface.
type Heap[T any] interface {
	// Push inserts one or more values into the heap while maintaining the heap property.
	//
	// Parameters:
	//   xs  The elements to be added to the heap.
	Push(xs ...T)

	// Pop removes and returns the element at the top of the heap (the one with
	// the highest priority).
	//
	// Returns:
	//   T     The element with the highest priority.
	//   bool  False if the heap is empty, true otherwise.
	Pop() (T, bool)

	// Peek returns the element at the top of the heap without removing it.
	//
	// Returns:
	//   T     The element with the highest priority.
	//   bool  False if the heap is empty, true otherwise.
	Peek() (T, bool)

	// Drain returns a destructive iterator that yields and removes elements
	// from the heap in priority order.
	//
	// The heap will be empty after the iterator completes its traversal.
	//
	// Returns:
	//   iter.Seq2[int, T]  A sequence of rank and elements from highest to lowest priority.
	Drain() iter.Seq2[int, T]

	pkg.Collection[T]
}
