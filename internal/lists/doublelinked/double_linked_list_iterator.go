package doublelinked

type doubleLinkedListIterator[T any] struct {
	current  *DoubleLinkedNode[T]
	reversed bool
}

func newDoubleLinkedListIterator[T any](
	first *DoubleLinkedNode[T], reversed bool) *doubleLinkedListIterator[T] {
	return &doubleLinkedListIterator[T]{
		current:  first,
		reversed: reversed}
}

func (iterator *doubleLinkedListIterator[T]) HasNext() bool {
	return iterator.current != nil
}

func (iterator *doubleLinkedListIterator[T]) Next() T {
	x := iterator.current.Value()
	if !iterator.reversed {
		iterator.current = iterator.current.Next()
		return x
	}

	iterator.current = iterator.current.Previous()
	return x
}
