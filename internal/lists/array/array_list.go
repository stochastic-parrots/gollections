package array

import (
	"fmt"
	"slices"
	"strings"

	"github.com/stochastic-parrots/gollections/internal/lists"
	"github.com/stochastic-parrots/gollections/pkg"
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

func (array ArrayList[T]) Get(index int) (T, error) {
	if index < 0 || index >= array.Length() {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, array.Length()-1)
	}

	return array.data[index], nil
}

func (array *ArrayList[T]) Set(index int, x T) error {
	if index < 0 || index >= array.Length() {
		return lists.NewIndexOutOfBoundError(index, array.Length()-1)
	}

	array.data[index] = x
	return nil
}

func (array *ArrayList[T]) Append(xs ...T) {
	for _, x := range xs {
		array.data = append(array.data, x)
		array.length++
	}
}

func (array ArrayList[T]) Iterator() pkg.Iterator[T] {
	return newArrayListIterator[T](array)
}

func (array *ArrayList[T]) Reverse() {
	slices.Reverse(array.data)
}

func (array ArrayList[T]) String() string {
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
