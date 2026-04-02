package collection_test

import (
	"encoding/json"
	"fmt"
	"iter"

	pkg "github.com/stochastic-parrots/gollections"
	"github.com/stochastic-parrots/gollections/internal/shared/collection"
)

var _ pkg.Collection[int] = &FakeCollection[int]{}

type FakeCollection[T any] struct {
	data []T
}

func (f *FakeCollection[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for idx, value := range f.data {
			if !yield(idx, value) {
				return
			}
		}
	}
}

func (f *FakeCollection[T]) IsEmpty() bool {
	return len(f.data) == 0
}

func (f *FakeCollection[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range f.data {
			if !yield(value) {
				return
			}
		}
	}
}

func (f *FakeCollection[T]) Length() int {
	return len(f.data)
}

func (f *FakeCollection[T]) Format(s fmt.State, verb rune) {
	collection.Format(s, verb, f, cap(f.data))
}

func (f *FakeCollection[T]) String() string {
	return fmt.Sprint(f)
}

func (f *FakeCollection[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.data)
}

func (f *FakeCollection[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.data)
}
