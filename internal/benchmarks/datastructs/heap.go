package datastructs

// Heap is the heap capability required by benchmark algorithms.
type Heap[T any] interface {
	Push(T)
	Pushes(...T)
	Length() int
	Pop() (T, bool)
	Peek() (T, bool)
	Replace(T) (T, bool)
}
