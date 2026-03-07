package binary_test

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/heaps"
	"github.com/stochastic-parrots/gollections/internal/heaps/binary"

	"container/heap"
	"math/rand"
)

// --- Boilerplate para a StdLib ---
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// --- Utils ---
var less = heaps.MinFunc[int]()

func BenchmarkBinaryHeap_Push(b *testing.B) {
	n := 10_000
	initialData := make([]int, n)
	ratios := []struct {
		name string
		val  float64
	}{
		{"x0.1", 0.1},
		{"x0.3", 0.3},
		{"x0.4", 0.4},
		{"x0.5", 0.5},
		{"x1", 1.0},
		{"x10", 10.0},
	}

	for _, r := range ratios {
		k := int(float64(n) * r.val)
		newItems := make([]int, k)

		b.Run(r.name, func(b *testing.B) {
			b.Run("Library=Gollections", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					h := binary.NewBinaryHeap(n+k, less)
					h.Push(initialData...)
					b.StartTimer()

					h.Push(newItems...)
				}
			})

			b.Run("Library=Container/Heap", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					h := &IntHeap{}
					*h = append(*h, initialData...)
					heap.Init(h)
					b.StartTimer()

					for _, item := range newItems {
						heap.Push(h, item)
					}
				}
			})
		})
	}
}

func BenchmarkBinaryHeap_Pop(b *testing.B) {
	const N = 100_000
	data := make([]int, N)
	for i := range N {
		data[i] = i
	}

	b.Run("Library=Gollections", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			h := binary.NewBinaryHeapFromSlice(append([]int{}, data...), less)
			b.StartTimer()

			for !h.IsEmpty() {
				h.Pop()
			}
		}
	})

	b.Run("Library=Container/Heap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			h := &IntHeap{}
			*h = append([]int{}, data...)
			heap.Init(h)
			b.StartTimer()

			for h.Len() > 0 {
				heap.Pop(h)
			}
		}
	})
}

func BenchmarkBinaryHeap_KLargest(b *testing.B) {
	const N = 100_000
	const K = 100

	data := make([]int, N)
	for i := range N {
		data[i] = rand.Int()
	}

	b.Run("Library=Gollections", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := binary.NewBinaryHeap(K+1, less)
			for _, x := range data {
				h.Push(x)
				if h.Length() > K {
					h.Pop()
				}
			}
		}
	})

	b.Run("Library=Container/Heap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			h := &IntHeap{}
			heap.Init(h)
			for _, x := range data {
				heap.Push(h, x)
				if h.Len() > K {
					heap.Pop(h)
				}
			}
		}
	})
}
