// Package constraint defines reusable generic type constraints.
package constraint

import "github.com/stochastic-parrots/gollections/internal/shared/constraint"

// Signed matches all built-in signed integer types.
type Signed = constraint.Signed

// Unsigned matches all built-in unsigned integer types.
type Unsigned = constraint.Unsigned

// Integer matches all built-in signed and unsigned integer types.
type Integer = constraint.Integer

// Float matches all built-in floating-point types.
type Float = constraint.Float

// Number matches all built-in integer and floating-point types.
type Number = constraint.Number

// Complex matches all built-in complex number types.
type Complex = constraint.Complex
