package doublelinked

import (
	"fmt"
	"iter"
	"strings"

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

func (list *DoubleLinkedList[T]) Length() int {
	return list.length
}

func (list *DoubleLinkedList[T]) IsEmpty() bool {
	return list.length == 0
}

func (list *DoubleLinkedList[T]) Get(index int) (T, error) {
	if index < 0 || index >= list.Length() {
		var zero T
		return zero, lists.NewIndexOutOfBoundError(index, list.Length()-1)
	}

	current := list.first
	for range index {
		if list.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	return current.Value(), nil
}

func (list *DoubleLinkedList[T]) Set(index int, x T) error {
	if index < 0 || index >= list.Length() {
		return lists.NewIndexOutOfBoundError(index, list.Length()-1)
	}

	current := list.first
	for range index {
		if list.reversed {
			current = current.Previous()
			continue
		}
		current = current.Next()
	}

	current.value = x
	return nil
}

func (list *DoubleLinkedList[T]) append(x T) {
	if list.IsEmpty() {
		list.first = NewDoubleLinkedNode(x)
		list.last = list.first
		list.length++
		return
	}

	if !list.reversed {
		new := NewDoubleLinkedNode(x)
		list.last.next = new
		new.previous = list.last
		list.last = new
		list.length++
		return
	}

	new := NewDoubleLinkedNode(x)
	list.last.previous = new
	new.next = list.last
	list.last = new
	list.length++
}

func (list *DoubleLinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		list.append(x)
	}
}

func (list *DoubleLinkedList[T]) Reverse() {
	temp := list.first
	list.first = list.last
	list.last = temp
	list.reversed = !list.reversed
}

func (list *DoubleLinkedList[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := list.first
		for current != nil {
			if !yield(current.value) {
				return
			}

			if !list.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

func (list *DoubleLinkedList[T]) Enumerate() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := list.first
		for index := 0; current != nil; index++ {
			if !yield(index, current.value) {
				return
			}

			if !list.reversed {
				current = current.next
			} else {
				current = current.previous
			}
		}
	}
}

func (list *DoubleLinkedList[T]) String() string {
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
