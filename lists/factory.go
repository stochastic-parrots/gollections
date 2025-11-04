package lists

import (
	lists "github.com/stochastic-parrots/gollections/internal/lists"
	linked "github.com/stochastic-parrots/gollections/internal/lists/linked"
)

type List[T any] = lists.List[T]

func NewLinkedList[T any]() List[T] {
	return linked.NewLinkedList[T]()
}
