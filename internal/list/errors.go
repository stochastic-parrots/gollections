package list

import "fmt"

// Operation identifies a list operation that uses an index.
type Operation uint8

const (
	// OperationGet identifies an indexed read.
	OperationGet Operation = iota + 1
	// OperationSet identifies an indexed update.
	OperationSet
	// OperationInsert identifies an indexed insertion.
	OperationInsert
	// OperationRemove identifies an indexed removal.
	OperationRemove
	// OperationReplace identifies an indexed replacement in a sorted list.
	OperationReplace
)

// IndexOutOfBoundsError reports an index outside the valid range [0, limit].
// Limit includes the end position for insertion errors.
type IndexOutOfBoundsError struct {
	index     int
	limit     int
	length    int
	operation Operation
}

// NewIndexOutOfBoundsError reports an invalid index for an operation.
// limit is the largest valid index; for insertion, it is the collection length.
func NewIndexOutOfBoundsError(operation Operation, index, limit int) *IndexOutOfBoundsError {
	length := limit + 1
	if operation == OperationInsert {
		length = limit
	}
	return &IndexOutOfBoundsError{
		index:     index,
		limit:     limit,
		length:    length,
		operation: operation,
	}
}

// Error reports the invalid index and the collection length.
func (e *IndexOutOfBoundsError) Error() string {
	operation := ""
	switch e.operation {
	case OperationGet:
		operation = "cannot get "
	case OperationSet:
		operation = "cannot set "
	case OperationInsert:
		operation = "cannot insert at "
	case OperationRemove:
		operation = "cannot remove "
	case OperationReplace:
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
