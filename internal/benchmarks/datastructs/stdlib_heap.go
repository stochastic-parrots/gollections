package datastructs

import "container/heap"

type stdLibIntHeap []int

func (h stdLibIntHeap) Len() int           { return len(h) }
func (h stdLibIntHeap) Length() int        { return len(h) }
func (h stdLibIntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h stdLibIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *stdLibIntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *stdLibIntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// StdLibHeap adapts container/heap to the benchmark Heap contract.
type StdLibHeap struct {
	stdheap *stdLibIntHeap
}

// Push inserts values into the standard-library heap baseline.
func (w *StdLibHeap) Push(xs ...int) {
	k := len(xs)
	if k == 0 {
		return
	}

	n := w.stdheap.Len()
	if n == 0 || (k > 64 && k > n) {
		*w.stdheap = append(*w.stdheap, xs...)
		heap.Init(w.stdheap)
		return
	}

	for _, x := range xs {
		heap.Push(w.stdheap, x)
	}
}

// Pop removes and returns the root of the standard-library heap baseline.
func (w *StdLibHeap) Pop() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return heap.Pop(w.stdheap).(int), true
}

// Peek returns the root of the standard-library heap baseline without removing it.
func (w *StdLibHeap) Peek() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return (*w.stdheap)[0], true
}

// Replace swaps the root and restores heap order.
func (w *StdLibHeap) Replace(v int) (int, bool) {
	if w.stdheap.Len() == 0 {
		w.Push(v)
		return -1, false
	}

	root := (*w.stdheap)[0]
	(*w.stdheap)[0] = v
	heap.Fix(w.stdheap, 0)

	return root, true
}

// Length returns the number of values in the heap baseline.
func (w *StdLibHeap) Length() int {
	return len(*w.stdheap)
}

// NewStdLibHeap creates an empty standard-library heap baseline.
func NewStdLibHeap(size int) *StdLibHeap {
	h := make(stdLibIntHeap, 0, size)
	heap.Init(&h)
	return &StdLibHeap{stdheap: &h}
}

// NewStdLibHeapCloneSlice creates a standard-library heap baseline from a cloned slice.
func NewStdLibHeapCloneSlice(data []int) *StdLibHeap {
	h := make(stdLibIntHeap, len(data))
	copy(h, data)
	heap.Init(&h)
	return &StdLibHeap{stdheap: &h}
}
