package disjointset

import "github.com/stochastic-parrots/gollections/internal/shared/constraint"

func flatInterval[T constraint.Integer](min, max T) int {
	if min > max {
		panic("disjointset: min must be <= max")
	}

	span := uint64(max) - uint64(min)
	maxInt := uint64(^uint(0) >> 1)

	if span >= maxInt {
		panic("disjointset: interval is too large")
	}

	return int(span) + 1
}

func flatIndex[T constraint.Integer](value, min T) int {
	return int(uint64(value) - uint64(min))
}

func flatValue[T constraint.Integer](index int, min T) T {
	return T(uint64(min) + uint64(index))
}
