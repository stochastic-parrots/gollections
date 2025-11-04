package array

import (
	"fmt"
	"slices"
	"strings"

	"github.com/stochastic-parrots/gollections"
)

type ArrayList[T any] struct {
	data   []T
	size   int
	length int
}

func NewArrayList[T any](size int) *ArrayList[T] {
	data := make([]T, 0, size)
	return &ArrayList[T]{data: data, size: size, length: 0}
}

func (array ArrayList[T]) Length() int {
	return array.length
}

func (array ArrayList[T]) IsEmpty() bool {
	return array.length == 0
}

func (array *ArrayList[T]) Append(xs ...T) {
	for _, x := range xs {
		array.data = append(array.data, x)
		array.length++
	}
}

func (array ArrayList[T]) Iterator() gollections.Iterator[T] {
	return newArrayListIterator(array)
}

func (array *ArrayList[T]) Reverse() {
	slices.Reverse(array.data)
}

func (array *ArrayList[T]) String() string {
	var builder strings.Builder
	builder.WriteRune('[')
	for i := range array.length {
		x := array.data[i]
		if i+1 < array.length {
			builder.WriteString(fmt.Sprintf("%v, ", x))
			continue
		}
		builder.WriteString(fmt.Sprintf("%v", x))
	}
	builder.WriteRune(']')
	return builder.String()
}
