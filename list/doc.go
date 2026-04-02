// Package list provides high-performance, generic linear collection implementations.
//
// The package offers a modern, type-safe suite of list structures, ranging from
// contiguous memory Array Lists to flexible Linked Lists. It leverages Go Generics (1.18+)
// to eliminate interface{} casting and the latest 'iter' package (1.23+) for
// high-performance, idiomatic data traversal.
//
// # Readonly Interface
//
// All lists implement the [Readonly] interface, providing a unified API:
//
//	type Readonly[T any] interface {
//		Get(idx int) (T, error)
//		Find(x T, cmp func(a, b T) int) (idx int, ok bool)
//		Contains(x T, cmp func(a, b T) int) bool
//		Backward() iter.Seq[T]
//		ToSlice() []T
//		pkg.Collection[T]
//		fmt.Stringer
//		json.Marshaler
//	}
//
// # List Interface
//
// All lists implement the [List] interface, providing a unified API:
//
//	type List[T any] interface {
//		Append(xs ...T)
//		Insert(idx int, x T) error
//		Set(idx int, x T) error
//		Remove(idx int) (T, error)
//		Reverse()
//		Clear()
//		Readonly[T]
//		json.Unmarshaler
//	}
//
// # Why this package?
//
//   - Performance: Optimized internal operations like O(1) Reverse for Double Linked Lists
//     and slice-backed efficiency for Array Lists.
//
//   - Modern Iteration: Full support for 'iter.Seq' and 'iter.Seq2', allowing you to
//     use range loops directly on your collections.
//
//   - Memory Management: Built with GC-friendly practices, including explicit
//     element zeroing to prevent memory leaks in generic types.
//
//   - Predictable API: Consistent method signatures across different list types
//     (Get, Set, Append, Length) to reduce the learning curve.
//
// # Core Concepts
//
// The package provides three main flavors of lists, each optimized for specific access patterns:
//
//   - [ArrayList]: Best for random access O(1) and memory locality.
//   - [LinkedList]: A doubly linked list optimized for efficient boundary insertions.
//     Supports O(1) reversal and bi-directional traversal.
//
// # Usage Example
//
//	l := list.NewLinked[string]()
//	l.Append("Go", "is", "awesome")
//
//	// Idiomatic iteration (Go 1.23+)
//	for i, val := range l.Enumerate() {
//	    fmt.Printf("Node %d: %s\n", i, val)
//	}
//
//	l.Reverse()
//	fmt.Println(l.String())
package list
