package datastructs

type Heap[T any] interface {
	Push(...T)
	Length() int
	Pop() (T, bool)
	Peek() (T, bool)
	Replace(T) (T, bool)
}
