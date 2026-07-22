package gollections

// Set defines the observation operations shared by collections whose values
// are unique according to an implementation-specific identity policy.
//
// Iteration order and value identity are defined by the concrete set. The
// index yielded by Enumerate is only the value's position in that iteration;
// unordered sets do not provide stable positional access.
type Set[T any] interface {
	Collection[T]

	// Contains returns true if the set contains a value with the same identity
	// as x according to the concrete implementation.
	Contains(x T) bool
}
