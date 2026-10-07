// Package deque provides generic double-ended queue implementations.
//
// A deque supports insertion and removal at both ends, making it useful as a
// queue, stack, sliding window, worklist, or small scheduling buffer. The package
// follows the same API shape as the rest of gollections: concrete implementations
// are exposed through direct constructors, while common traversal uses Go's iter
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
//	}
//
// # Implementations
//
// The package exposes two implementations:
//
//   - [ArrayDeque]: A circular-array deque with good memory locality and
//     amortized O(1) single-value operations at both ends.
//   - [LinkedDeque]: A doubly linked deque with O(1) single-value operations
//     at both ends and no backing-array moves.
//
// Construct an [ArrayDeque] with [NewArray], [ArrayFrom], [ArrayClone], or
// [ArrayFromSeq]. Construct a [LinkedDeque] with [NewLinked], [LinkedFrom], or
// [LinkedFromSeq]. The zero values of both concrete types are ready for use.
//
// [ArrayFrom] takes ownership of the provided slice in front-to-back order;
// the caller must stop using the slice and every alias of its backing array.
// [ArrayClone] makes an independent shallow copy. [LinkedFrom] copies values
// into new nodes without modifying or retaining the source. Both FromSeq
// constructors consume their iterators once in front-to-back order.
//
// # JSON
//
// ArrayDeque and LinkedDeque marshal and unmarshal as arrays from front to back.
// The Deque and Readonly interfaces do not require JSON or formatting methods.
// Unmarshal replaces the current contents only after the complete JSON array
// is decoded successfully.
//
// All traverses a deque from front to back. Enumerate uses the same order
// and assigns consecutive indexes starting at zero.
package deque
