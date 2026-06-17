package algorithms

import (
	"math"

	"github.com/stochastic-parrots/gollections/constraint"
)

func infinity[T constraint.Number]() T {
	var zero T
	switch any(zero).(type) {
	case float64:
		return any(math.Inf(1)).(T)
	case float32:
		return any(float32(math.Inf(1))).(T)
	case int:
		return any(math.MaxInt).(T)
	case int64:
		return any(int64(math.MaxInt64)).(T)
	case int32:
		return any(int32(math.MaxInt32)).(T)
	case int16:
		return any(int16(math.MaxInt16)).(T)
	case int8:
		return any(int8(math.MaxInt8)).(T)
	case uint:
		return any(uint(math.MaxUint)).(T)
	case uint64:
		return any(uint64(math.MaxUint64)).(T)
	case uint32:
		return any(uint32(math.MaxUint32)).(T)
	case uint16:
		return any(uint16(math.MaxUint16)).(T)
	case uint8:
		return any(uint8(math.MaxUint8)).(T)
	case uintptr:
		return any(^uintptr(0)).(T)
	default:
		return zero
	}
}
