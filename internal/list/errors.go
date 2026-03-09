package list

import "fmt"

// IndexOutOfBoundError represents an error where the provided index is outside
// the valid range [0, limit].
//
// This struct implements the 'error' interface and should be used by list methods
// (like Get or Set) to return context-rich errors.
type IndexOutOfBoundError struct {
	// index is the invalid index that the user attempted to access.
	index int

	// limit is the maximum valid index allowed in the collection (Length - 1).
	limit int
}

// NewIndexOutOfBoundError creates and returns a new instance of IndexOutOfBoundError.
//
// This function should be used internally by the library to generate an
// out-of-bounds error with contextual information.
//
// Parameters:
//
//	index  The invalid index that was requested.
//	limit  The largest valid index allowed in the collection (Length - 1).
func NewIndexOutOfBoundError(index, limit int) *IndexOutOfBoundError {
	return &IndexOutOfBoundError{index, limit}
}

// Error implements the 'error' interface and formats the error message,
// including the invalid index and the maximum allowed limit.
func (e *IndexOutOfBoundError) Error() string {
	return fmt.Sprintf("index %d is out of bounds; maximum valid index is %d", e.index, e.limit)
}

// ErrIndexOutOfBound is the "sentinel" error used for type checking.
//
// API consumers should use errors.Is(err, lists.ErrIndexOutOfBound)
// to check if the returned error is an IndexOutOfBoundError.
//
// To extract the 'index' and 'limit' fields, use errors.As.
var ErrIndexOutOfBound = &IndexOutOfBoundError{}
