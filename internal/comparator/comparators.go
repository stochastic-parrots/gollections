package comparator

import "cmp"

func Min[T cmp.Ordered]() func(a, b T) bool { return func(a, b T) bool { return a < b } }

func Max[T cmp.Ordered]() func(a, b T) bool { return func(a, b T) bool { return a > b } }
