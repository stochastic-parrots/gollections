package suites

import (
	"container/heap"

	internal "github.com/stochastic-parrots/gollections/internal/heap"
)

// --- Boilerplate StdLib Heap ---
type StdLibPureHeap []int

func (h StdLibPureHeap) Len() int           { return len(h) }
func (h StdLibPureHeap) Length() int        { return len(h) }
func (h StdLibPureHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h StdLibPureHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *StdLibPureHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *StdLibPureHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type StdLibHeap struct {
	stdheap *StdLibPureHeap
}

func (w *StdLibHeap) Push_(x int) {
	w.stdheap.Push(x)
}

func (w *StdLibHeap) Push(xs ...int) {
	for _, x := range xs {
		w.stdheap.Push(x)
	}
}

func (w *StdLibHeap) Pop() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return w.stdheap.Pop().(int), true
}

func (w *StdLibHeap) Peek() (int, bool) {
	if w.stdheap.Len() == 0 {
		return -1, false
	}
	return (*w.stdheap)[0], true
}

func (w *StdLibHeap) Length() int {
	return w.stdheap.Len()
}

func NewStdLibHeap(data []int) *StdLibHeap {
	h := &StdLibPureHeap{}
	*h = append(*h, data...)
	heap.Init(h)
	return &StdLibHeap{stdheap: h}
}

// --- Boilerplate StdLib Priority Map ---
type StdLibEntry[K comparable, V any] struct {
	key      K
	priority V
}

type StdPriorityMapHeap[K comparable, V any] struct {
	data    []StdLibEntry[K, V]
	less    func(a, b V) bool
	indices map[K]int
}

func (h StdPriorityMapHeap[K, V]) Len() int { return len(h.data) }
func (h StdPriorityMapHeap[K, V]) Less(i, j int) bool {
	return h.less(h.data[i].priority, h.data[j].priority)
}
func (h StdPriorityMapHeap[K, V]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.indices[h.data[i].key] = i
	h.indices[h.data[j].key] = j
}
func (h *StdPriorityMapHeap[K, V]) Push(x any) {
	it := x.(StdLibEntry[K, V])
	h.indices[it.key] = len(h.data)
	h.data = append(h.data, it)
}
func (h *StdPriorityMapHeap[K, V]) Pop() any {
	old := h.data
	n := len(old)
	it := old[n-1]
	h.data = old[0 : n-1]
	delete(h.indices, it.key)
	return it
}

type StdPriorityMap[K comparable, V any] struct {
	heap *StdPriorityMapHeap[K, V]
}

func NewStdPriorityMap[K comparable, V any](size int, less func(a, b V) bool) *StdPriorityMap[K, V] {
	return &StdPriorityMap[K, V]{
		heap: &StdPriorityMapHeap[K, V]{
			less:    less,
			data:    make([]StdLibEntry[K, V], 0, size),
			indices: make(map[K]int, size),
		},
	}
}

func (m *StdPriorityMap[K, V]) Set(key K, priority V) {
	if idx, ok := m.heap.indices[key]; ok {
		m.heap.data[idx].priority = priority
		heap.Fix(m.heap, idx)
		return
	}
	heap.Push(m.heap, StdLibEntry[K, V]{key, priority})
}

func (m *StdPriorityMap[K, V]) Pop() (K, V, bool) {
	if m.heap.Len() == 0 {
		var zk K
		var zv V
		return zk, zv, false
	}

	item := heap.Pop(m.heap).(StdLibEntry[K, V])
	return item.key, item.priority, true
}

func (m *StdPriorityMap[K, V]) Length() int {
	return len(m.heap.data)
}

// --- Utils ---
var less = internal.MinFunc[int]()

type Implementation[T any] struct {
	name    string
	factory func() T
}

type Implementations[T any] []Implementation[T]
