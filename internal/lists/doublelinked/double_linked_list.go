package doublelinked

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
	"github.com/stochastic-parrots/gollections/internal/lists"
)

type DoubleLinkedList[T any] struct {
	first, last *DoubleLinkedNode[T]
	length      int
	reversed    bool
}

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		first:    nil,
		last:     nil,
		length:   0,
		reversed: false,
	}
}

func (l *DoubleLinkedList[T]) Length() int {
	return l.length
}

func (l *DoubleLinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *DoubleLinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.Length() {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		if l.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	return current.Value(), nil
}

func (l *DoubleLinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= l.Length() {
		return lists.NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		if l.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	current.value = x
	return nil
}

func (l *DoubleLinkedList[T]) append(x T) {
	if l.IsEmpty() {
		l.first = NewDoubleLinkedNode(x)
		l.last = l.first
		l.length++
		return
	}

	if !l.reversed {
		new := NewDoubleLinkedNode(x)
		l.last.next = new
		new.previous = l.last
		l.last = new
		l.length++
		return
	}

	new := NewDoubleLinkedNode(x)
	l.last.previous = new
	new.next = l.last
	l.last = new
	l.length++
}

func (l *DoubleLinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		l.append(x)
	}
}

func (l *DoubleLinkedList[T]) Reverse() {
	temp := l.first
	l.first = l.last
	l.last = temp
	l.reversed = !l.reversed
}

func (l *DoubleLinkedList[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := l.first
		for current != nil {
			if !yield(current.value) {
				return
			}

			if !l.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

func (l *DoubleLinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := l.first
		for index := 0; current != nil; index++ {
			if !yield(index, current.value) {
				return
			}

			if !l.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

func (l *DoubleLinkedList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, l.Length())
}

func (l *DoubleLinkedList[T]) String() string {
	return fmt.Sprint(l)
}
