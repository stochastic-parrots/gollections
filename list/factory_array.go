package list

import (
	"iter"

	"github.com/stochastic-parrots/gollections/internal/list"
)

// ArrayList is a slice-backed [List]. Its zero value is ready for use.
type ArrayList[T any] = list.ArrayList[T]

var _ List[int] = &list.ArrayList[int]{}

// NewArray creates an empty slice-backed list. Choose it for O(1) indexed
// access and writes, amortized O(1) append, and memory-local traversal.
// Insertion or removal at an arbitrary position is O(N).
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Get(idx)              O(1)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(1)
//	Remove(idx)           O(N)
//	Reverse()             O(N)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(capacity).
func NewArray[T any](capacity int) *ArrayList[T] { return list.NewArrayList[T](capacity) }

// ArrayFrom transfers ownership of data and its backing array; the caller must stop using the slice and every alias afterward. Use [ArrayClone] to preserve the source.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Get(idx)              O(1)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(1)
//	Remove(idx)           O(N)
//	Reverse()             O(N)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(1).
func ArrayFrom[T any](data []T) *ArrayList[T] { return list.NewArrayListFromSlice(data) }

// ArrayClone makes an independent shallow copy of data without modifying or retaining its backing array.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Get(idx)              O(1)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(1)
//	Remove(idx)           O(N)
//	Reverse()             O(N)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(len(data)).
func ArrayClone[T any](data []T) *ArrayList[T] { return list.NewArrayListCloneSlice(data) }

// ArrayFromSeq consumes seq once into new storage in iteration order.
//
// Performance Summary (Time Complexity):
//
//	Operation             Time Complexity
//	-------------------   ---------------
//	IsEmpty()             O(1)
//	Len()              O(1)
//	All()                 O(N)
//	Enumerate()           O(N)
//	Get(idx)              O(1)
//	Find(x, cmp)          O(N)
//	Contains(x, cmp)      O(N)
//	Backward()            O(N)
//	ToSlice()             O(N)
//	Append(x)             O(1) Amortized
//	Appends(xs... T)      O(len(xs)) Amortized
//	Insert(idx, x)        O(N)
//	Set(idx, x)           O(1)
//	Remove(idx)           O(N)
//	Reverse()             O(N)
//	Clear()               O(N)
//	UnmarshalJSON(data)   O(N + len(data))
//	MarshalJSON()         O(N)
//	String()              O(1)
//
// Complexity: O(N), where N is the number of values yielded by seq.
func ArrayFromSeq[T any](seq iter.Seq[T]) *ArrayList[T] { return list.NewArrayListFromSeq(seq) }
