// Package comparator provides types and utilities for comparing values.
package comparator

import "cmp"

// Ordering is a predicate function that returns true if a should
// come before b in a specific order (e.g., a < b).
//
// This matches the signature of cmp.Less and is the primary
// type used in PriorityMap and Heap factories.
type Ordering[T any] func(a, b T) bool

// --- Helpers for Ordering (bool) ---

// Min returns a less-than predicate for ordered types.
func Min[T cmp.Ordered]() Ordering[T] {
	return func(a, b T) bool { return a < b }
}

// Max returns a greater-than predicate for ordered types.
func Max[T cmp.Ordered]() Ordering[T] {
	return func(a, b T) bool { return a > b }
}
