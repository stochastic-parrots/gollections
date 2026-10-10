package disjointset

import "github.com/stochastic-parrots/gollections/internal/shared/constraint"

func flatInterval[T constraint.Integer](minimum, maximum T) int {
	if minimum > maximum {
		panic("disjointset: min must be <= max")
	}

	span := uint64(maximum) - uint64(minimum)
	maxInt := uint64(^uint(0) >> 1)

	if span >= maxInt {
		panic("disjointset: interval is too large")
	}

	return int(span) + 1
}

func flatIndex[T constraint.Integer](value, minimum T) int {
	return int(uint64(value) - uint64(minimum))
}

func flatValue[T constraint.Integer](index int, minimum T) T {
	return T(uint64(minimum) + uint64(index))
}
