package doublelinked

type DoubleLinkedNode[T any] struct {
	previous *DoubleLinkedNode[T]
	value    T
	next     *DoubleLinkedNode[T]
}

func NewDoubleLinkedNode[T any](value T) *DoubleLinkedNode[T] {
	return &DoubleLinkedNode[T]{
		previous: nil,
		value:    value,
		next:     nil,
	}
}

func (node *DoubleLinkedNode[T]) Unlink() {
	previous := node.previous
	next := node.next

	if previous != nil {
		previous.next = next
	}
	if next != nil {
		next.previous = previous
	}

	node.previous = nil
	node.next = nil
}

func (node DoubleLinkedNode[T]) Value() T {
	return node.value
}

func (node DoubleLinkedNode[T]) Next() *DoubleLinkedNode[T] {
	return node.next
}

func (node *DoubleLinkedNode[T]) Append(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.next = new
	new.previous = node
	return new
}

func (node DoubleLinkedNode[T]) HasNext() bool {
	return node.next != nil
}

func (node DoubleLinkedNode[T]) Previous() *DoubleLinkedNode[T] {
	return node.previous
}

func (node *DoubleLinkedNode[T]) PreAppend(x T) *DoubleLinkedNode[T] {
	new := NewDoubleLinkedNode(x)
	node.previous = new
	new.next = node
	return new
}

func (node DoubleLinkedNode[T]) HasPrevious() bool {
	return node.previous != nil
}
