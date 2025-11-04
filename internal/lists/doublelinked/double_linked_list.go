package doublelinked

import (
	"fmt"
	"strings"

	"github.com/stochastic-parrots/gollections"
)

type DoubleLinkedList[T any] struct {
	First    *DoubleLinkedNode[T]
	Last     *DoubleLinkedNode[T]
	Len      int
	reversed bool
}

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		First:    nil,
		Last:     nil,
		Len:      0,
		reversed: false,
	}
}

func (list DoubleLinkedList[T]) Length() int {
	return list.Len
}

func (list DoubleLinkedList[T]) IsEmpty() bool {
	return list.Len == 0
}

func (list *DoubleLinkedList[T]) append(x T) {
	if list.IsEmpty() {
		list.First = NewDoubleLinkedNode(x)
		list.Last = list.First
		list.Len++
		return
	}

	if !list.reversed {
		new := NewDoubleLinkedNode(x)
		list.Last.next = new
		new.previous = list.Last
		list.Last = new
		list.Len++
		return
	}

	new := NewDoubleLinkedNode(x)
	list.Last.previous = new
	new.next = list.Last
	list.Last = new
	list.Len++
}

func (list *DoubleLinkedList[T]) Append(xs ...T) {
	for _, x := range xs {
		list.append(x)
	}
}

func (list *DoubleLinkedList[T]) Reverse() {
	temp := list.First
	list.First = list.Last
	list.Last = temp
	list.reversed = !list.reversed
}

func (list DoubleLinkedList[T]) Iterator() gollections.Iterator[T] {
	return newDoubleLinkedListIterator(list.First, list.reversed)
}

func (list DoubleLinkedList[T]) String() string {
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
