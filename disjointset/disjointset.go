package disjointset

import (
	"iter"

	"github.com/stochastic-parrots/gollections"
	"github.com/stochastic-parrots/gollections/constraint"
)

// Readonly describes observation of a fixed, inclusive integer range partitioned
// into disjoint sets. Find and Connected may compress paths internally, so even
// calls through a readonly view require synchronization with concurrent access.
type Readonly[T constraint.Integer] interface {
	gollections.Collection[T]

	// Find returns the current representative of x. The representative may
	// change after a union. It returns the zero value and false when x is
	// outside the range.
	Find(x T) (representative T, ok bool)

	// Connected reports whether a and b belong to the same set. The second
	// result is false when either value is outside the range.
	Connected(a, b T) (connected, ok bool)

	// Disjoints returns the current number of disjoint sets.
	Disjoints() int
}

// DisjointSet partitions a fixed, inclusive integer range into sets that can
// be merged but not individually split.
type DisjointSet[T constraint.Integer] interface {
	Readonly[T]

	// Union merges the sets containing a and b. It returns true only when the
	// sets were distinct and both values belong to the range.
	Union(a, b T) bool

	// Reset restores every value in the range to its own set.
	Reset()
}

// RankDisjointSet adds a rank query to the common disjoint-set operations.
// Rank is a union heuristic and does not report the size of a set.
type RankDisjointSet[T constraint.Integer] interface {
	DisjointSet[T]

	// Rank returns the rank of the set containing x, or zero and false when x
	// is outside the range.
	Rank(x T) (rank int, ok bool)
}

// SizeDisjointSet adds a set-size query to the common disjoint-set operations.
type SizeDisjointSet[T constraint.Integer] interface {
	DisjointSet[T]

	// Size returns the number of values in the set containing x, or zero and
	// false when x is outside the range.
	Size(x T) (size int, ok bool)
}

// AsReadonly returns a [Readonly] view of the provided [DisjointSet]. The view
// observes the same underlying partition, but does not expose Union or Reset.
// It is neither a snapshot nor synchronization. Nil input returns nil.
func AsReadonly[T constraint.Integer](set DisjointSet[T]) *readonly[T] {
	if set == nil {
		return nil
	}
	return &readonly[T]{inner: set}
}

type readonly[T constraint.Integer] struct {
	inner DisjointSet[T]
}

func (w readonly[T]) Find(x T) (T, bool) { return w.inner.Find(x) }

func (w readonly[T]) Connected(a, b T) (bool, bool) { return w.inner.Connected(a, b) }

func (w readonly[T]) Disjoints() int { return w.inner.Disjoints() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Len() int { return w.inner.Len() }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

var _ Readonly[int] = (*readonly[int])(nil)
