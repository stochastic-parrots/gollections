// Package gollections provides a suite of high-performance, generic data structures for Go.
//
// The library leverages Go generics for type safety and standard iterators for
// idiomatic data traversal. Instead of a "one-size-fits-all" approach,
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
// Root-level interfaces are intentionally read-only capability contracts. They
// describe how callers can inspect or iterate over a structure without changing
// its state. Mutating operations such as Append, Push, Set, Remove, Pop, or
// Clear are exposed only by the structure-specific interfaces in each
// subpackage, where their semantics are precise and unambiguous.
//
// # Subpackages
//
// The library is organized into specialized subpackages. Refer to each package
// documentation for detailed time complexity tables:
//
//   - [github.com/stochastic-parrots/gollections/list]:
//     Indexed sequences like ArrayList and LinkedList.
//
//   - [github.com/stochastic-parrots/gollections/deque]:
//     Double-ended queues backed by circular arrays or linked nodes.
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
//   - Read-only Roots: Shared root interfaces expose observation and iteration;
//     mutation belongs to the concrete collection family.
//
//   - Performance: Focused on O(1) and O(log N) operations where possible.
//
//   - Idiomatic Go: Full support for range over iterators.
//
// # Documentation Characteristics
//
// Each data structure includes detailed performance documentation in its subpackage.
// Refer to the factory functions (for example, list.NewArray) for time complexity
// tables covering all operations.
package gollections
