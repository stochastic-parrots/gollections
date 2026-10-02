package set

import "iter"

// Source provides sized iteration for internal set algorithms.
type Source[T any] interface {
	All() iter.Seq[T]
	Length() int
}
