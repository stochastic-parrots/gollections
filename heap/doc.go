// Package heap provides high-performance, generic heap implementations.
//
// The package offers a modern alternative to the standard library's container/heap,
// leveraging Go Generics (1.18+) for type safety and the new iterators (1.23+)
// for idiomatic data traversal.
//
// # Heap Interface
//
// All heaps implement the Heap interface, providing a unified API:
//
//	type Heap[T any] interface {
//		Push(xs ...T)
//		Pop() (T, bool)
//		Peek() (T, bool)
//		Drain() iter.Seq2[int, T]
//		Replace(x T) (T, bool)
//		Collection[T]
//	}
//
// # Why this package?
//
//   - Type Safety: No more interface{} casting; work directly with your types.
//
//   - Go Idiomatic: Integrates with the 'iter' package for seamless range loops.
//
//   - Versatile: Supports Min-Heaps, Max-Heaps, and custom priority logic.
//
//   - Memory Efficient: Provides in-place constructors with sufix 'From'
//     and safe copies constructors with sufix 'Clone'.
//
//   - Predictable API: Consistent method signatures across different heap types
//     (Push, Pop, Peek, Drain) to reduce the learning curve.
//
// # Core Concepts
//
// The [Heap] interface extends the base [gollections.Collection], adding priority-specific
// operations like Push, Pop, and Peek.
//
// One of the standout features is the Drain method:
//
//	h := heap.MinBinaryFrom([]int{3, 1, 2})
//	for i, val := range h.Drain() {
//	    fmt.Printf("Rank %d: %v\n", i, val)
//	}
//
// After this loop, h is empty.
package heap
