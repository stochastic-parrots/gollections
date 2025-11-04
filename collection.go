package gollections

import "fmt"

type Collection[T any] interface {
	IsEmpty() bool
	Length() int
	Reverse()
	Iterator() Iterator[T]
	fmt.Stringer
}
