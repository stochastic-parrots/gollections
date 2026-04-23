package node

// LinkedNode represents an individual singly linked node.
// It stores a generic value and a pointer to the subsequent node in the sequence.
type LinkedNode[T any] struct {
	Value T
	Next  *LinkedNode[T]
}

// NewLinkedNode creates a new LinkedNode initialized with the provided value.
//
// Complexity: O(1).
func NewLinkedNode[T any](value T) *LinkedNode[T] {
	return &LinkedNode[T]{Value: value, Next: nil}
}

// HasNext returns true if there is a node following the current one.
//
// Complexity: O(1).
func (node LinkedNode[T]) HasNext() bool {
	return node.Next != nil
}

// Append creates a new node with the given value and attaches it
// immediately after the current node.
//
// Complexity: O(1).
func (node *LinkedNode[T]) Append(x T) *LinkedNode[T] {
	new := NewLinkedNode(x)
	node.Next = new
	return new
}
