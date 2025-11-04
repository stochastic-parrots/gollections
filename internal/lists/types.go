package lists

import "github.com/stochastic-parrots/gollections"

type List[T any] interface {
	Append(...T)
	gollections.Collection[T]
}
