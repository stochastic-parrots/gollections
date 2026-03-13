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

type StdLibHeap struct {
	stdheap *stdLibIntHeap
}

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

func (w *StdLibHeap) Pop() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return heap.Pop(w.stdheap).(int), true
}

func (w *StdLibHeap) Peek() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return (*w.stdheap)[0], true
}

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

func (w *StdLibHeap) Length() int {
	return len(*w.stdheap)
}

func NewStdLibHeap(size int) *StdLibHeap {
	h := make(stdLibIntHeap, 0, size)
	heap.Init(&h)
	return &StdLibHeap{stdheap: &h}
}

func NewStdLibHeapCloneSlice(data []int) *StdLibHeap {
	h := make(stdLibIntHeap, len(data))
	copy(h, data)
	heap.Init(&h)
	return &StdLibHeap{stdheap: &h}
}
