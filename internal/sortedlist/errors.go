package sortedlist

import "errors"

// ErrOrderViolation reports a replacement that would break sorted order.
var ErrOrderViolation = errors.New("replacement violates sorted order")
