// Package set provides mutable generic collections of unique values.
//
// # Implementations
//
// [HashSet] stores comparable values directly as Go map keys. Choose it when Go
// equality defines value identity. Its zero value is ready for use. Select it
// with [HashSetOf].
//
// [KeyedHashSet] stores arbitrary values and derives a comparable identity key
// through a caller-provided function. Choose it for values that are not
// comparable, or when only part of a value defines membership. The first value
// inserted for a key remains its representative. Its zero value is invalid;
// select it with [HashSetBy].
//
// A KeyedHashSet identity must remain stable while a value belongs to the set. For
// example, when pointers are stored and the identity function reads a field,
// callers must not mutate that field until the value is removed.
//
// Hash identity must also be reflexive. Go floating-point NaN values are not
// equal to themselves, so they cannot be found after insertion into HashSet and
// must not be returned directly as KeyedHashSet keys. Use HashSetBy with a
// canonical comparable representation when NaN values need set membership.
// When HashSet uses an interface type, every dynamic value supplied as a key
// must be comparable, matching the requirements of a native Go map.
//
// # Construction
//
// HashSetOf and HashSetBy return reusable typed factories. New creates an empty
// set, From copies values from a slice into map storage, and FromSeq collects
// an iterator. From and FromSeq discard duplicates while preserving the first
// value encountered for each identity.
//
// # Set Algebra
//
// HashSet and KeyedHashSet implement [Algebra]. Clone, Union, Intersection,
// Difference, and SymmetricDifference return independent sets without modifying
// their receivers. Calling a receiver operation without additional operands
// returns a clone. The factories expose the same operations for arbitrary
// [Source] values and return an empty set when called without sources. Source
// requires only sized iteration; operands do not need to implement the full
// [gollections.Collection] contract.
//
// Concrete sets in this package satisfy [gollections.Collection] through
// pointer methods. Pass *HashSet or *KeyedHashSet operands, as returned by their
// factories; the corresponding struct values do not implement the operand
// interface, and initialized set structs must not be copied. Concrete pointer
// operands also let the implementations traverse map storage directly instead
// of allocating iterator adapters. Operations that normalize derived
// identities may still allocate a compact membership map because the receiver's
// key function must be reapplied.
//
// Every source is interpreted under the receiver's identity policy, or under
// the factory's policy for factory operations. Duplicate identities within an
// arbitrary collection are therefore treated as one membership. For
// KeyedHashSet, Union preserves the first representative encountered;
// Intersection and Difference preserve representatives from the receiver, or
// from the first source when called through a factory. SymmetricDifference uses
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
// New and that capacity before calling UnionWith avoids map growth.
//
// Equal, subset, superset, and disjoint relations also apply the selected
// identity policy to both operands. Consequently, receiver relations are
// directional when comparing sets constructed with different identity
// policies; use a factory relation when both operands should be interpreted by
// one explicitly selected policy.
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
// identity policy. Decode into []T and pass the values to the same HashSetOf or
// HashSetBy factory.
//
// Sets are not safe for concurrent use. Callers must synchronize access when
// any goroutine may mutate a shared set.
package set
