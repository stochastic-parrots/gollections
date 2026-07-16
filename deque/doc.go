// Package deque provides generic double-ended queue implementations.
//
// A deque supports insertion and removal at both ends, making it useful as a
// queue, stack, sliding window, worklist, or small scheduling buffer. The package
// follows the same API shape as the rest of gollections: concrete implementations
// are hidden behind public factories, while common traversal uses Go's iter
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
// A readonly view observes the same underlying deque. It restricts mutation
// through that interface but is not a snapshot or concurrency mechanism.
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
//
// Select a strategy with [Array] or [Linked], then construct a concrete deque
// through the returned factory. Array factories provide New, From, Clone, and
// FromSeq. Linked factories provide New, From, and FromSeq because linked
// construction always copies values into new nodes.
// The zero values of [ArrayDeque] and [LinkedDeque] are ready for use.
//
// [ArrayFactory.From] takes ownership of the provided slice and interprets its
// values in front-to-back order. Clone preserves the source slice.
// [LinkedFactory.From] never modifies or retains its source.
//
// # JSON
//
// Deques marshal and unmarshal as arrays from front to back. Unmarshal replaces
// the current contents only after the complete JSON array is decoded
// successfully.
//
// All traverses a deque from front to back. Enumerate uses the same order
// and assigns consecutive indexes starting at zero.
package deque
