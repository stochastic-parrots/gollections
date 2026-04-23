package node

// DoubleLinkedNode represents a single doubly linked node.
// It contains a value of type T and pointers to the previous and next nodes
// in the sequence, allowing for bidirectional navigation.
type DoubleLinkedNode[T any] struct {
	Previous *DoubleLinkedNode[T]
	Value    T
	Next     *DoubleLinkedNode[T]
}

// NewDoubleLinkedNode creates an empty DoubleLinkedNode ready for use.
//
// Complexity: O(1).
func NewDoubleLinkedNode[T any](value T) *DoubleLinkedNode[T] {
	return &DoubleLinkedNode[T]{
		Previous: nil,
		Value:    value,
		Next:     nil,
	}
}

// Unlink disconnects the node from its neighbors.
//
// Complexity: O(1).
func (node *DoubleLinkedNode[T]) Unlink() {
	previous := node.Previous
	next := node.Next

	if previous != nil {
		previous.Next = next
	}
	if next != nil {
		next.Previous = previous
	}

	node.Previous = nil
	node.Next = nil
}

// Append creates a new node with the given value and attaches it
// immediately after the current node.
//
// Complexity: O(1).
func (node *DoubleLinkedNode[T]) Append(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.Next = new
	new.Previous = node
	return new
}

// HasNext returns true if there is a node following the current one.
//
// Complexity: O(1).
func (node DoubleLinkedNode[T]) HasNext() bool {
	return node.Next != nil
}

// PreAppend creates a new node with the given value and attaches it
// immediately before the current node.
//
// Complexity: O(1).
func (node *DoubleLinkedNode[T]) Prepend(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.Previous = new
	new.Next = node
	return new
}

// HasPrevious returns true if there is a node preceding the current one.
//
// Complexity: O(1).
func (node DoubleLinkedNode[T]) HasPrevious() bool {
	return node.Previous != nil
}
