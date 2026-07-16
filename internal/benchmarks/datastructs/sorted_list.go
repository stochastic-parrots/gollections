package datastructs

// SortedList is the sorted-list capability required by benchmark algorithms.
type SortedList interface {
	Add(int)
	Adds(...int)
	Contains(int) bool
	Remove(int) bool
	Get(idx int) (int, error)
	LowerBound(int) int
	Clear()
}
