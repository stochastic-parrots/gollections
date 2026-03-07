// Package gollections provides a collection of generic data structures for Go.
//
// Gollections offers efficient, type-safe implementations of common collections
// such as lists and heaps, leveraging Go's generics (available in Go 1.18+).
// All data structures are fully generic, allowing you to work with any type
// without sacrificing performance or type safety.
//
// # Lists
//
// The lists subpackage provides three implementations optimized for different use cases:
//
//   - ArrayList: Slice-based list with O(1) random access and amortized O(1) append.
//     Best for scenarios requiring frequent random access.
//
//   - LinkedList: Singly linked list with O(N) random access but efficient insertion/deletion
//     at known positions. Useful for sequential traversal and memory efficiency.
//
//   - DoubleLinkedList: Doubly linked list supporting efficient traversal in both directions.
//     Best when you need bidirectional iteration or frequent modifications near the end.
//
// All lists support common operations like Append, Get, Set, Reverse, Sort, and iteration.
// They implement the Collection interface, ensuring consistent API across implementations.
//
// The heaps subpackage provides min-heaps and max-heaps with support for custom comparators.
// Heaps are useful for priority queues, sorting, and selection problems.
//
//   - Min Heap: Efficiently retrieves the smallest element.
//   - Max Heap: Efficiently retrieves the largest element.
//   - Custom Heap: Define your own comparison function for complex sorting logic.
//
// Heaps can be created from scratch or initialized from existing slices. When initialized
// from a slice, the heap operates on that slice directly for minimal memory overhead.
//
// # Design Principles.
//
//   - Type Safety: Full generic support ensures compile-time type checking.
//   - Performance: Implementations are optimized for their intended use cases.
//   - Simplicity: Clean, intuitive APIs follow Go conventions (e.g., iter.Seq iterators).
//   - Flexibility: Support for custom comparators and multiple implementations.
//
// # Collection Interface
//
// Most collections implement the Collection interface, providing a unified API:
//
//	type Collection[T any] interface {
//		IsEmpty() bool
//		Length() int
//		Iterator() iter.Seq[T]
//		fmt.Stringer
//	}
//
// Additional operations like Reverse and Sort are available on specific implementations
// (e.g., all lists support these). Use type assertions if needed to access implementation-specific methods.
//
// # Iterators
//
// All collections support Go's iter.Seq[T] for efficient, range-based iteration.
// This allows seamless integration with the range keyword (Go 1.22+):
//
//	for val := range list.Iterator() {
//		fmt.Println(val)
//	}
//
// # Performance Characteristics
//
// Each data structure includes detailed performance documentation in its subpackage.
// Refer to the factory functions (e.g., lists.NewArrayList) for Time Complexity tables
// covering all operations.
package gollections
