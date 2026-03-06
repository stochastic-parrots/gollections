package linked

import (
	"fmt"
	"iter"
	"strings"

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

func (list LinkedList[T]) Length() int {
	return list.length
}

func (list LinkedList[T]) IsEmpty() bool {
	return list.length == 0
}

func (list LinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= list.Length() {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, list.Length()-1)
	}

	current := list.first
	for range index {
		current = current.Next()
	}

	return current.Value(), nil
}

func (list LinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= list.Length() {
		return lists.NewIndexOutOfBoundError(index, list.Length()-1)
	}

	current := list.first
	for range index {
		current = current.Next()
	}

	current.value = x
	return nil
}

func (list *LinkedList[T]) append(x T) {
	new := NewNode(x)

	if list.IsEmpty() {
		list.first = new
		list.last = new
		list.length++
		return
	} else {
		list.last.next = new
		list.last = new
	}
	list.length++
}

func (list *LinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		list.append(x)
	}
}

func (list *LinkedList[T]) Reverse() {
	if list.Length() <= 1 {
		return
	}

	var previous *Node[T]
	current := list.first

	for current != nil {
		next := current.next
		current.next = previous
		previous = current
		current = next
	}

	tmp := list.first
	list.first = previous
	list.last = tmp
}

func (list LinkedList[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		for current := list.first; current != nil; current = current.next {
			if !yield(current.value) {
				return
			}
		}
	}
}

func (list LinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for current, index := list.first, 0; current != nil; current = current.next {
			if !yield(index, current.value) {
				return
			}
			index++
		}
	}
}

func (list LinkedList[T]) String() string {
	if list.IsEmpty() {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteRune('[')

	for i, val := range list.Enumerate() {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, "%v", val)
	}

	sb.WriteRune(']')
	return sb.String()
}
