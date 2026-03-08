package pkg

import (
	"fmt"
	"iter"
)

// Collection defines the basic operations for a generic data structure.
type Collection[T any] interface {
	// IsEmpty Checks if the collection is empty.
	IsEmpty() bool

	// Length Returns the collection length (number of elements).
	Length() int

	// Iterator Returns an iterator for the collection.
	Iterator() iter.Seq[T]

	// Iterator Returns an indexed iterator for the collection.
	Enumerate() iter.Seq2[int, T]

	fmt.Stringer

	fmt.Formatter
}
