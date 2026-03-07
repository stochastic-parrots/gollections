package formatters_test

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
	"github.com/stochastic-parrots/gollections/pkg"
)

var _ pkg.Collection[int] = &FakeCollection[int]{}

type FakeCollection[T any] struct {
	data []T
}

func (fake *FakeCollection[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range fake.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

func (fake *FakeCollection[T]) IsEmpty() bool {
	return len(fake.data) == 0
}

func (fake *FakeCollection[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range fake.data {
			if !yield(value) {
				return
			}
		}
	}
}

func (fake *FakeCollection[T]) Length() int {
	return len(fake.data)
}

func (fake *FakeCollection[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, fake, cap(fake.data))
}

func (fake *FakeCollection[T]) String() string {
	return fmt.Sprint(fake)
}
