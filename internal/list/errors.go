package list

import "fmt"

// IndexOperation identifies an indexed list operation.
type IndexOperation uint8

const (
	// IndexOperationGet identifies an indexed read.
	IndexOperationGet IndexOperation = iota + 1
	// IndexOperationSet identifies an indexed update.
	IndexOperationSet
	// IndexOperationInsert identifies an indexed insertion.
	IndexOperationInsert
	// IndexOperationRemove identifies an indexed removal.
	IndexOperationRemove
	// IndexOperationReplace identifies an indexed replacement in a sorted list.
	IndexOperationReplace
)

// IndexOutOfBoundsError reports an index outside the valid range [0, limit].
// Limit includes the end position for insertion errors.
type IndexOutOfBoundsError struct {
	index     int
	limit     int
	length    int
	operation IndexOperation
}

// NewIndexOutOfBoundsError reports an invalid index and the largest valid index.
// For indexed reads, updates, and removals, limit is the collection length minus one.
func NewIndexOutOfBoundsError(index, limit int) *IndexOutOfBoundsError {
	return &IndexOutOfBoundsError{index: index, limit: limit, length: limit + 1}
}

// NewIndexOperationOutOfBoundsError reports an invalid index for a specific operation.
// limit is the largest valid index; for insertion, it is the collection length.
func NewIndexOperationOutOfBoundsError(operation IndexOperation, index, limit int) *IndexOutOfBoundsError {
	return &IndexOutOfBoundsError{
		index:     index,
		limit:     limit,
		length:    limit + 1,
		operation: operation,
	}
}

// Error reports the invalid index and the collection length.
func (e *IndexOutOfBoundsError) Error() string {
	operation := ""
	switch e.operation {
	case IndexOperationGet:
		operation = "cannot get "
	case IndexOperationSet:
		operation = "cannot set "
	case IndexOperationInsert:
		operation = "cannot insert at "
	case IndexOperationRemove:
		operation = "cannot remove "
	case IndexOperationReplace:
		operation = "cannot replace "
	}
	return fmt.Sprintf("%sindex %d out of range for length %d", operation, e.index, e.length)
}

// Is reports whether the target is an IndexOutOfBoundsError.
func (e *IndexOutOfBoundsError) Is(target error) bool {
	_, ok := target.(*IndexOutOfBoundsError)
	return ok
}

// Index returns the invalid index that was requested.
func (e *IndexOutOfBoundsError) Index() int {
	return e.index
}

// Limit returns the largest valid index, including the end position for insertion.
func (e *IndexOutOfBoundsError) Limit() int {
	return e.limit
}

// ErrIndexOutOfBounds is the "sentinel" error used for type checking.
//
// API consumers should use errors.Is(err, list.ErrIndexOutOfBounds)
// to check if the returned error is an IndexOutOfBoundsError.
//
// To extract the index and limit values, use errors.As and the Index/Limit methods.
var ErrIndexOutOfBounds = &IndexOutOfBoundsError{}
