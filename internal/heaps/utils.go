package heaps

import "cmp"

func MinFunc[T cmp.Ordered]() func(a, b T) bool { return func(a, b T) bool { return a < b } }

func MaxFunc[T cmp.Ordered]() func(a, b T) bool { return func(a, b T) bool { return a > b } }
