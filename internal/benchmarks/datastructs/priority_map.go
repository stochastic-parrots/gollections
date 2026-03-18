package datastructs

type PriorityMap[K comparable, V any] interface {
	Set(key K, priority V)
	Pop() (key K, priority V, ok bool)
	Length() int
	Clear()
}
