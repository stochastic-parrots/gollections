package algorithms

type Heap[T any] interface {
	Push(...T)
	Length() int
	Pop() (T, bool)
	Peek() (T, bool)
}
