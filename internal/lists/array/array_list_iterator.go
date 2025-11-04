package array

import "github.com/stochastic-parrots/gollections"

type arrayListIterator[T any] struct {
	index  int
	length int
	data   *[]T
}

func newArrayListIterator[T any](array ArrayList[T]) gollections.Iterator[T] {
	return &arrayListIterator[T]{
		index:  0,
		length: array.length,
		data:   &array.data,
	}
}

func (i *arrayListIterator[T]) HasNext() bool {
	return i.index < i.length
}

func (i *arrayListIterator[T]) Next() T {
	current := (*i.data)[i.index]
	i.index++
	return current
}
