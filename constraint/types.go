// Package constraint defines reusable generic type constraints.
package constraint

// Signed matches all built-in signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned matches all built-in unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer matches all built-in signed and unsigned integer types.
type Integer interface {
	Signed | Unsigned
}

// Float matches all built-in floating-point types.
type Float interface {
	~float32 | ~float64
}

// Number matches all built-in integer and floating-point types.
type Number interface {
	Integer | Float
}

// Complex matches all built-in complex number types.
type Complex interface {
	~complex64 | ~complex128
}
