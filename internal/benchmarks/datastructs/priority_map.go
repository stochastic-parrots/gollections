package datastructs

// PriorityMap is the keyed priority queue capability required by graph benchmarks.
type PriorityMap[K comparable, V any] interface {
	Set(key K, priority V)
	Pop() (key K, priority V, ok bool)
	Length() int
	Clear()
}
