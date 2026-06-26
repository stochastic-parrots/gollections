package sortedlist

import (
	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

// IndexOutOfBoundError reports an index outside the valid sorted-list range.
type IndexOutOfBoundError = list.IndexOutOfBoundError

// ErrIndexOutOfBound identifies sorted-list index bounds errors.
var ErrIndexOutOfBound = list.ErrIndexOutOfBound

// ErrOrderViolation identifies replacements that would break sorted order.
var ErrOrderViolation = sortedlist.ErrOrderViolation
