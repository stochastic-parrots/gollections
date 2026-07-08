package datastructs

// SortedList is the sorted-list capability required by benchmark algorithms.
type SortedList interface {
	Add(...int)
	Remove(int) bool
	Contains(int) bool
	Get(idx int) (int, error)
	LowerBound(int) int
	Clear()
}
