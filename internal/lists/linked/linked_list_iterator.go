package linked

import "github.com/stochastic-parrots/gollections"

type linkedListIterator[T any] struct {
	current *Node[T]
}

func newIterator[T any](first *Node[T]) gollections.Iterator[T] {
	return &linkedListIterator[T]{current: first}
}

func (i *linkedListIterator[T]) HasNext() bool {
	return i.current != nil
}

func (i *linkedListIterator[T]) Next() T {
	x := i.current.Value()
	i.current = i.current.Next()
	return x
}
