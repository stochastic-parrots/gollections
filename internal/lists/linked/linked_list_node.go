package linked

type Node[T any] struct {
	value T
	next  *Node[T]
}

func NewNode[T any](value T) *Node[T] {
	return &Node[T]{value: value, next: nil}
}

func (node Node[T]) Value() T {
	return node.value
}

func (node Node[T]) Next() *Node[T] {
	return node.next
}

func (node Node[T]) HasNext() bool {
	return node.next != nil
}

func (node *Node[T]) Append(x T) *Node[T] {
	new := NewNode(x)
	node.next = new
	return new
}
