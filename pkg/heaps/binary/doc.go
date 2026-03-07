// Package binary provides a high-performance Binary Heap implementation of the heaps.Heap interface.
//
// A Binary Heap is a complete binary tree stored in a contiguous slice, offering
// an excellent balance between pointer-free memory efficiency and logarithmic performance.
//
// # Performance Characteristics
//
// This implementation uses a zero-indexed array representation where for any node at index i:
//   - Left child:  2i + 1
//   - Right child: 2i + 2
//   - Parent:      (i - 1) / 2
//
// By leveraging Go Generics, this package eliminates the "interface tax" (boxing/unboxing)
// found in traditional container/heap implementations.
//
// This provides three main advantages:
//  1. Zero Allocations: Operations on primitive types do not escape to the heap.
//  2. Inlining: The compiler can inline comparison logic, reducing function call overhead.
//  3. Cache Locality: Contiguous memory access is optimized for modern CPU caches.
//
// # Heapify and Construction
//
// When using [NewBinaryHeapFromSlice], this package utilizes Floyd's Heapify algorithm,
// which converts an unordered slice into a valid heap in O(N) time. This is
// significantly faster than N individual Push operations (O(N log N)).
//
// # Memory Management
//
//   - NewBinaryHeapFromSlice: Reuses the provided slice directly. It is highly memory
//     efficient but destructive to the original element order.
//   - NewBinaryHeapCloneSlice: Allocates a new slice, preserving the original data.
//
// # Usage Example
//
//	h := binary.Min[int](10)
//	h.Push(5, 3, 8)
//	val, ok := h.Pop() // returns 3, true
package binary
