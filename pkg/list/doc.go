// Package list provides high-performance, generic linear collection implementations.
//
// The package offers a modern, type-safe suite of list structures, ranging from
// contiguous memory ArrayLists to flexible Linked Lists. It leverages Go Generics (1.18+)
// to eliminate interface{} casting and the latest 'iter' package (1.23+) for
// high-performance, idiomatic data traversal.
//
// # Why this package?
//
//   - Performance: Optimized internal operations like O(1) Reverse for DoubleLinkedLists
//     and slice-backed efficiency for ArrayLists.
//   - Modern Iteration: Full support for 'iter.Seq' and 'iter.Seq2', allowing you to
//     use range loops directly on your collections.
//   - Memory Management: Built with GC-friendly practices, including explicit
//     element zeroing to prevent memory leaks in generic types.
//   - Predictable API: Consistent method signatures across different list types
//     (Get, Set, Append, Length) to reduce the learning curve.
//
// # Core Concepts
//
// The package provides three main flavors of lists, each optimized for specific access patterns:
//
//   - [array.ArrayList]: Best for random access (O(1)) and memory locality.
//   - [linked.LinkedList]: A classic singly linked list for efficient head insertions.
//   - [doublelinked.DoubleLinkedList]: Supports O(1) reversal and bi-directional traversal.
//
// # Usage Example
//
//	l := linked.NewLinkedList[string]()
//	l.Append("Go", "is", "awesome")
//
//	// Idiomatic iteration (Go 1.23+)
//	for i, val := range l.Enumerate() {
//	    fmt.Printf("Node %d: %s\n", i, val)
//	}
//
//	l.Reverse()
//	fmt.Println(l.String())
//
// # Formatting
//
// All collections in this package implement the fmt.Formatter interface, providing
// clean, customizable output when using the fmt package's verbs (%v, %+v, %#v).
package list
