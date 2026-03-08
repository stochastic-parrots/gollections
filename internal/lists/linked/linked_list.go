package linked

import (
	"fmt"
	"iter"

	"github.com/stochastic-parrots/gollections/internal/formatters"
	"github.com/stochastic-parrots/gollections/internal/lists"
)

type LinkedList[T any] struct {
	first, last *Node[T]
	length      int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		first:  nil,
		last:   nil,
		length: 0,
	}
}

func (l *LinkedList[T]) Length() int {
	return l.length
}

func (l *LinkedList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= l.Length() {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		current = current.Next()
	}

	return current.Value(), nil
}

func (l *LinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= l.Length() {
		return lists.NewIndexOutOfBoundError(index, l.Length()-1)
	}

	current := l.first
	for range index {
		current = current.Next()
	}

	current.value = x
	return nil
}

func (l *LinkedList[T]) append(x T) {
	new := NewNode(x)

	if l.IsEmpty() {
		l.first = new
		l.last = new
		l.length++
		return
	} else {
		l.last.next = new
		l.last = new
	}
	l.length++
}

func (l *LinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		l.append(x)
	}
}

func (l *LinkedList[T]) Reverse() {
	if l.Length() <= 1 {
		return
	}

	var previous *Node[T]
	current := l.first

	for current != nil {
		next := current.next
		current.next = previous
		previous = current
		current = next
	}

	tmp := l.first
	l.first = previous
	l.last = tmp
}

func (l *LinkedList[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for current := l.first; current != nil; current = current.next {
			if !yield(current.value) {
				return
			}
		}
	}
}

func (l *LinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for current, index := l.first, 0; current != nil; current = current.next {
			if !yield(index, current.value) {
				return
			}
			index++
		}
	}
}

func (l *LinkedList[T]) Format(s fmt.State, verb rune) {
	formatters.Format(s, verb, l, l.Length())
}

func (l *LinkedList[T]) String() string {
	return fmt.Sprint(l)
}
