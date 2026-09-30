package sortedlist

import (
	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// IndexOutOfBoundsError reports an index outside the valid sorted-list range.
type IndexOutOfBoundsError = list.IndexOutOfBoundsError

// ErrIndexOutOfBounds identifies sorted-list index bounds errors.
var ErrIndexOutOfBounds = list.ErrIndexOutOfBounds

// ErrOrderViolation identifies replacements that would break sorted order.
var ErrOrderViolation = sortedlist.ErrOrderViolation
