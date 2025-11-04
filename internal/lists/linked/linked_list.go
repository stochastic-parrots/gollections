package linked

import (
	"fmt"
	"strings"

	"github.com/stochastic-parrots/gollections"
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

func (list *LinkedList[T]) append(x T) {
	new := NewNode(x)

	if list.IsEmpty() {
		list.first = new
		list.last = new
		list.length++
		return
	}

	current := list.first
	for current.HasNext() {
		current = current.Next()
	}

	current.next = new
	list.last = new
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

func (list LinkedList[T]) Iterator() gollections.Iterator[T] {
	return newIterator(list.first)
}

func (list LinkedList[T]) String() string {
	var sb strings.Builder
	var index int

	sb.WriteRune('[')
	for it := list.Iterator(); it.HasNext(); {
		if index+1 < list.Length() {
			sb.WriteString(fmt.Sprintf("%v, ", it.Next()))
			index++
			continue
		}

		sb.WriteString(fmt.Sprintf("%v", it.Next()))
		index++
	}

	sb.WriteRune(']')
	return sb.String()
}
