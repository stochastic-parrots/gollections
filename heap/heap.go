package heap

import (
	"encoding/json"
	"fmt"
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
)

// Heap defines a priority queue whose top element is selected by a comparator.
//
// The comparator supplied by the concrete constructor defines what "highest
// priority" means, so the same interface supports min-heaps, max-heaps, and
// custom ordering.
type Heap[T any] interface {
	// Push inserts one or more values into the heap while maintaining the heap property.
	Push(xs ...T)

	// Pop removes and returns the element at the top of the heap (the one with
	// the highest priority).
	//
	// It returns the zero value of T and false if the heap is empty.
	Pop() (T, bool)

	// Peek returns the element at the top of the heap without removing it.
	//
	// It returns the zero value of T and false if the heap is empty.
	Peek() (T, bool)

	// Drain returns a destructive iterator that yields and removes elements
	// from the heap in priority order.
	//
	// If the iterator is fully consumed, the heap will be empty. If iteration
	// stops early, only the yielded elements are removed.
	//
	// The yielded index is the zero-based rank of each removed element.
	Drain() iter.Seq2[int, T]

	// Replace replaces the root element with the given item and rebalances the heap.
	//
	// It returns the previous root and true. If the heap is empty, it returns
	// the zero value of T and false without inserting x.
	Replace(x T) (T, bool)

	// Clear removes all elements from the heap, resetting it to an empty state.
	Clear()

	pkg.Collection[T]

	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
}
