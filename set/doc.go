// Package set provides mutable generic collections of unique values.
//
// # Implementations
//
// [HashSet] stores comparable values directly as Go map keys. Choose it when Go
// equality defines value identity. Its zero value is ready for use. Use
// [NewHashSet] to choose an initial capacity or [HashSetFrom] to populate it.
//
// [KeyedHashSet] stores arbitrary values and derives a comparable identity key
// through a caller-provided function. Choose it for values that are not
// comparable, or when only part of a value defines membership. The first value
// inserted for a key remains its representative. Its zero value is invalid;
// construct it with [NewKeyedHashSet].
//
// A KeyedHashSet identity must remain stable while a value belongs to the set. For
// example, when pointers are stored and the identity function reads a field,
// callers must not mutate that field until the value is removed.
//
// Hash identity must also be reflexive. Go floating-point NaN values are not
// equal to themselves, so they cannot be found after insertion into HashSet and
// must not be returned directly as KeyedHashSet keys. Use NewKeyedHashSet with
// a canonical comparable representation when NaN values need set membership.
// When HashSet uses an interface type, every dynamic value supplied as a key
// must be comparable, matching the requirements of a native Go map.
//
// # Construction
//
// NewHashSet creates an empty HashSet, while HashSetFrom and HashSetFromSeq
// build one from a slice or iterator. NewKeyedHashSet creates an empty
// KeyedHashSet, while KeyedHashSetFrom and KeyedHashSetFromSeq build one from a
// slice or iterator. KeyedHashSetFromMap shares a map as its backing storage;
// KeyedHashSetCloneMap copies it. For both map constructors, callers are
// responsible for ensuring each map key matches the identity of its value.
// Slice and iterator constructors discard duplicates while preserving the
// first value encountered for each identity.
//
// # Set Algebra
//
// HashSet and KeyedHashSet implement [Algebra]. Clone, Union, Intersection,
// Difference, and SymmetricDifference return independent sets without modifying
// their receivers. Calling a receiver operation without additional operands
// returns a clone. Set algebra accepts arbitrary [Source] values, which require
// only sized iteration and need not implement the full
// [gollections.Collection] contract.
// [SourceFromSlice] returns a [SliceSource] that wraps a slice as an operand
// without copying it.
//
// Concrete sets in this package satisfy [gollections.Collection] through
// pointer methods. Pass *HashSet or *KeyedHashSet operands, as returned by their
// constructors; the corresponding struct values do not implement the operand
// interface, and initialized set structs must not be copied. Concrete pointer
// operands also let the implementations traverse map storage directly instead
// of allocating iterator adapters. Operations that normalize derived
// identities may still allocate a compact membership map because the receiver's
// key function must be reapplied.
//
// Every source is interpreted under the receiver's identity policy. Duplicate
// identities within an arbitrary collection are therefore treated as one
// membership. For
// KeyedHashSet, Union preserves the first representative encountered;
// Intersection and Difference preserve representatives from the receiver.
// SymmetricDifference uses
// the representative from the last source that toggles an identity into the
// result. When that source has unspecified iteration order and multiple values
// collapse to the same selected identity, the chosen representative is
// unspecified.
//
// [InPlaceAlgebra] provides UnionWith, IntersectWith, DifferenceWith, and
// SymmetricDifferenceWith when reusing an existing set is preferable to
// allocating an independent result. These methods are destructive, are safe
// when an operand aliases the receiver, and return the number of receiver
// memberships changed. With no operands they leave the receiver unchanged. If
// the expected final size of a union is known, creating an empty receiver with
// [NewHashSet] or [NewKeyedHashSet] and that capacity before calling UnionWith
// avoids map growth.
//
// Equal, IsSubset, IsProperSubset, IsSuperset, IsProperSuperset, and IsDisjoint
// are package functions over [Readonly] sets that use direct Go equality.
// Their EqualBy, IsSubsetBy, IsProperSubsetBy, IsSupersetBy,
// IsProperSupersetBy, and IsDisjointBy counterparts accept an identity
// function and apply it to both operands. KeyedHashSet also exposes relations
// as methods that use the receiver's identity function. Receiver relations are
// directional when comparing sets constructed with different identity
// policies. The By functions require a non-nil keyOf that returns reflexive
// comparable keys.
//
// # Views, Iteration, And JSON
//
// [Set] exposes mutation, while [Readonly] contains observation and traversal
// operations. [AsReadonly] prevents callers from recovering the mutable set
// through a type assertion. The view observes the same underlying set; it is
// not a snapshot and does not make concurrent mutation safe.
//
// All and Enumerate iterate in unspecified order. Enumerate assigns consecutive
// positions starting at zero for that traversal only; those positions are not
// stable set indexes. Iterators are lazy and callers must not mutate the set
// while iteration is running.
//
// Sets marshal as JSON arrays in unspecified iteration order. They do not
// implement json.Unmarshaler because a JSON array does not encode the selected
// identity policy. Decode into []T and pass the values to HashSetFrom or
// KeyedHashSetFrom with the same identity function.
//
// Sets are not safe for concurrent use. Callers must synchronize access when
// any goroutine may mutate a shared set.
package set
