package set

import (
	"encoding/json"
	"iter"

	"github.com/stochastic-parrots/gollections"
	"github.com/stochastic-parrots/gollections/internal/set"
)

// Source defines the sized iteration required by set algebra operands.
type Source[T any] = set.Source[T]

// Readonly defines a non-mutable view of a set.
type Readonly[T any] interface {
	gollections.Collection[T]
	json.Marshaler

	// Contains returns true if the set contains a value with the same identity
	// as x according to the concrete implementation.
	Contains(x T) bool
}

// Set defines a mutable collection of values that are unique according to its
// concrete implementation's identity policy.
type Set[T any] interface {
	Readonly[T]

	// Add inserts x if the set does not already contain a value with the same
	// identity. Existing values are preserved. It returns true if x was added.
	Add(x T) (added bool)

	// Adds inserts every value in xs that does not already exist in the set.
	// When multiple values have the same identity, the first one is preserved.
	// It returns the number of values added.
	Adds(xs ...T) (added int)

	// Remove deletes the value with the same identity as x.
	// It returns true if a value was removed.
	Remove(x T) bool

	// Removes deletes every value in xs whose identity belongs to the set.
	// It returns the number of values removed.
	Removes(xs ...T) (removed int)

	// Clear removes all values and leaves the set ready for reuse with the same
	// identity policy.
	Clear()
}

// Algebra defines pure set operations whose results preserve a concrete set
// type and its identity policy. Implementations do not mutate their operands.
//
// Every operand is interpreted as a set under the receiver's identity policy,
// so duplicate identities in arbitrary collections are normalized. Relations
// are directional when operands use different identity policies. HashSet and
// KeyedHashSet satisfy operand interfaces through pointers; callers should pass
// the pointers returned by their factories rather than struct values.
type Algebra[T any, S Readonly[T]] interface {
	// Clone returns an independent set containing the receiver's values and
	// preserving its identity policy.
	Clone() S

	// Union returns a new set containing identities present in the receiver or
	// any other operand.
	Union(others ...Source[T]) S

	// Intersection returns a new set containing identities present in the
	// receiver and every other operand.
	Intersection(others ...Source[T]) S

	// Difference returns a new set containing receiver identities absent from
	// every other operand.
	Difference(others ...Source[T]) S

	// SymmetricDifference returns a new set containing identities present in an
	// odd number of the receiver and other operands.
	SymmetricDifference(others ...Source[T]) S

	// Equal returns true if the receiver and other contain the same identities.
	Equal(other Source[T]) bool

	// IsSubset returns true if every receiver identity is present in other.
	IsSubset(other Source[T]) bool

	// IsProperSubset returns true if the receiver is a subset of other and the
	// two operands are not equal.
	IsProperSubset(other Source[T]) bool

	// IsSuperset returns true if every identity in other is present in the
	// receiver.
	IsSuperset(other Source[T]) bool

	// IsProperSuperset returns true if the receiver is a superset of other and
	// the two operands are not equal.
	IsProperSuperset(other Source[T]) bool

	// IsDisjoint returns true if the receiver and other share no identity.
	IsDisjoint(other Source[T]) bool
}

// InPlaceAlgebra defines destructive set operations. Every operand is
// interpreted under the receiver's identity policy. Calling a method without
// operands leaves the receiver unchanged, and an operand may alias the
// receiver. HashSet and KeyedHashSet operands must be passed as pointers.
type InPlaceAlgebra[T any] interface {
	// UnionWith adds identities from every operand and returns how many were
	// added.
	UnionWith(others ...Source[T]) int

	// IntersectWith removes receiver identities absent from any operand and
	// returns how many were removed.
	IntersectWith(others ...Source[T]) int

	// DifferenceWith removes identities present in any operand and returns how
	// many were removed.
	DifferenceWith(others ...Source[T]) int

	// SymmetricDifferenceWith toggles identities present in an odd number of
	// operands and returns how many receiver memberships changed.
	SymmetricDifferenceWith(others ...Source[T]) int
}

// AsReadonly returns a [Readonly] view of the provided [Set].
//
// The returned view observes the same underlying set while preventing type
// assertion back to the mutable interface. It is not a snapshot and does not
// provide synchronization.
func AsReadonly[T any](set Set[T]) *readonly[T] {
	if set == nil {
		return nil
	}
	return &readonly[T]{inner: set}
}

type readonly[T any] struct {
	inner Set[T]
}

func (w readonly[T]) Contains(x T) bool { return w.inner.Contains(x) }

func (w readonly[T]) All() iter.Seq[T] { return w.inner.All() }

func (w readonly[T]) Enumerate() iter.Seq2[int, T] { return w.inner.Enumerate() }

func (w readonly[T]) IsEmpty() bool { return w.inner.IsEmpty() }

func (w readonly[T]) Length() int { return w.inner.Length() }

func (w readonly[T]) MarshalJSON() ([]byte, error) { return w.inner.MarshalJSON() }

var _ Readonly[any] = (*readonly[any])(nil)
