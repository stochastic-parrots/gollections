package array

import (
	"fmt"
	"iter"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/formatters"
	"github.com/stochastic-parrots/gollections/internal/lists"
)

type ArrayList[T any] struct {
	data []T
}

func NewArrayList[T any](size int) *ArrayList[T] {
	data := make([]T, 0, size)
	return &ArrayList[T]{data: data}
}

func (l *ArrayList[T]) Length() int {
	return len(l.data)
}

func (l *ArrayList[T]) IsEmpty() bool {
	return len(l.data) == 0
}

func (l *ArrayList[T]) Get(index int) (T, error) {
	if index < 0 || index >= len(l.data) {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	return l.data[index], nil
}

func (l *ArrayList[T]) Set(index int, x T) error {
	if index < 0 || index >= len(l.data) {
		return lists.NewIndexOutOfBoundError(index, len(l.data)-1)
	}

	l.data[index] = x
	return nil
}

func (l *ArrayList[T]) Append(xs ...T) {
	l.data = append(l.data, xs...)
}

func (l *ArrayList[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range l.data {
			if !yield(value) {
				return
			}
		}
	}
}

func (l *ArrayList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range l.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

func (l *ArrayList[T]) Reverse() {
	slices.Reverse(l.data)
}

func (l *ArrayList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, cap(l.data))
}

func (l *ArrayList[T]) String() string {
	return fmt.Sprint(l)
}
