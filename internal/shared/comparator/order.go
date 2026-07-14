// Package comparator provides shared ordering functions for internal data
// structures.
package comparator

import "cmp"

// Min returns a less-than predicate for ordered types.
func Min[T cmp.Ordered]() func(T, T) bool {
	return cmp.Less[T]
}

// Max returns a greater-than predicate for ordered types.
func Max[T cmp.Ordered]() func(T, T) bool {
	return func(a, b T) bool { return cmp.Less(b, a) }
}
