// Package gollections provides a suite of high-performance, generic data structures for Go.
//
// The library leverages Go Generics (1.18+) for type safety and the standard iterators
// (1.23+) for idiomatic data traversal. Instead of a "one-size-fits-all" approach,
// gollections provides specialized implementations optimized for specific memory
// and performance profiles.
//
// # Core Interfaces
//
// Most data structures in this module implement one of the two base interfaces,
// ensuring a predictable API across different implementations:
//
//   - [Collection]: Foundation for linear structures like lists and heaps.
//   - [Map]: Base operations for key-value based structures.
//
// # Subpackages
//
// The library is organized into specialized subpackages. Refer to each package
// documentation for detailed time complexity tables:
//
//   - [github.com/stochastic-parrots/gollections/list]:
//     Indexed sequences like ArrayList and LinkedList.
//
//   - [github.com/stochastic-parrots/gollections/heap]:
//     Priority-based ordering (Min/Max Heaps).
//
//   - [github.com/stochastic-parrots/gollections/prioritymap]:
//     A hybrid structure combining Map lookups with Heap ordering.
//
// # Design Principles
//
//   - Type Safety: Full generic support ensures compile-time type checking.
//
//   - Performance: Focused on O(1) and O(log N) operations where possible.
//
//   - Idiomatic Go: Full support for 'range' over iterators (Go 1.23+).
//
// # Documentation Characteristics
//
// Each data structure includes detailed performance documentation in its subpackage.
// Refer to the factory functions (e.g., list.NewArrayList) for Time Complexity
// tables covering all operations.
package gollections
