// Package deque provides generic double-ended queue implementations.
//
// A deque supports insertion and removal at both ends, making it useful as a
// queue, stack, sliding window, worklist, or small scheduling buffer. The package
// follows the same API shape as the rest of gollections: concrete implementations
// are hidden behind public constructors, while common traversal uses Go's iter
// package.
//
// # Readonly Interface
//
// All deques implement the [Readonly] interface:
//
//	type Readonly[T any] interface {
//		Front() (T, bool)
//		Back() (T, bool)
//		ToSlice() []T
//		gollections.Collection[T]
//		fmt.Stringer
//		json.Marshaler
//	}
//
// # Deque Interface
//
// Mutable deques implement the [Deque] interface:
//
//	type Deque[T any] interface {
//		Prepend(xs ...T)
//		Append(xs ...T)
//		Shift() (T, bool)
//		Pop() (T, bool)
//		Clear()
//		Readonly[T]
//		json.Unmarshaler
//	}
//
// # Implementations
//
// The package exposes two implementations:
//
//   - [ArrayDeque]: A circular-array deque with good memory locality and
//     amortized O(1) operations at both ends.
//   - [LinkedDeque]: A doubly linked deque with stable O(1) operations at both
//     ends and no backing-array moves.
package deque
