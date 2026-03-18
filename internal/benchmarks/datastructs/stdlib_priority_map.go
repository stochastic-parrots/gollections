package datastructs

import "container/heap"

type stdLibEntry[K comparable, V any] struct {
	key      K
	priority V
}

type stdPriorityMapHeap[K comparable, V any] struct {
	data    []stdLibEntry[K, V]
	less    func(a, b V) bool
	indexes map[K]int
}

func (h stdPriorityMapHeap[K, V]) Len() int {
	return len(h.data)
}

func (h stdPriorityMapHeap[K, V]) Less(i, j int) bool {
	return h.less(h.data[i].priority, h.data[j].priority)
}

func (h stdPriorityMapHeap[K, V]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.indexes[h.data[i].key] = i
	h.indexes[h.data[j].key] = j
}

func (h *stdPriorityMapHeap[K, V]) Push(x any) {
	it := x.(stdLibEntry[K, V])
	h.indexes[it.key] = len(h.data)
	h.data = append(h.data, it)
}

func (h *stdPriorityMapHeap[K, V]) Pop() any {
	old := h.data
	n := len(old)
	it := old[n-1]
	h.data = old[0 : n-1]
	delete(h.indexes, it.key)
	return it
}

type StdPriorityMap[K comparable, V any] struct {
	heap *stdPriorityMapHeap[K, V]
}

func NewStdPriorityMap[K comparable, V any](size int, less func(a, b V) bool) *StdPriorityMap[K, V] {
	return &StdPriorityMap[K, V]{
		heap: &stdPriorityMapHeap[K, V]{
			less:    less,
			data:    make([]stdLibEntry[K, V], 0, size),
			indexes: make(map[K]int, size),
		},
	}
}

func (m *StdPriorityMap[K, V]) Set(key K, priority V) {
	if idx, ok := m.heap.indexes[key]; ok {
		m.heap.data[idx].priority = priority
		heap.Fix(m.heap, idx)
		return
	}
	heap.Push(m.heap, stdLibEntry[K, V]{key, priority})
}

func (m *StdPriorityMap[K, V]) Pop() (K, V, bool) {
	if m.heap.Len() == 0 {
		var zk K
		var zv V
		return zk, zv, false
	}

	item := heap.Pop(m.heap).(stdLibEntry[K, V])
	return item.key, item.priority, true
}

func (m *StdPriorityMap[K, V]) Length() int {
	return len(m.heap.data)
}

func (m *StdPriorityMap[K, V]) Clear() {
	clear(m.heap.indexes)
	clear(m.heap.data)
}
