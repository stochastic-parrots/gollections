// Package heaps provides high-performance, generic binary heap implementations.
//
// The package offers a modern alternative to the standard library's container/heap,
// leveraging Go Generics (1.18+) for type safety and the new iterators (1.23+)
// for idiomatic data traversal.
//
// # Why this package?
//
//   - Type Safety: No more interface{} casting; work directly with your types.
//   - Go Idiomatic: Integrates with the 'iter' package for seamless range loops.
//   - Versatile: Supports Min-Heaps, Max-Heaps, and custom priority logic.
//   - Memory Efficient: Provides both in-place construction (FromSlice) and
//     safe copies (CloneSlice).
//
// # Core Concepts
//
// The [Heap] interface extends the base [pkg.Collection], adding priority-specific
// operations like Push, Pop, and Peek.
//
// One of the standout features is the Drain method:
//
//	h := heaps.MinFromSlice([]int{3, 1, 2})
//	for i, val := range h.Drain() {
//	    fmt.Printf("Rank %d: %v\n", i, val)
//	}
//
// After this loop, h is empty.
package heaps
