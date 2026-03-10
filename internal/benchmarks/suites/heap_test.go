package suites

import (
	"cmp"
	"math/rand/v2"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/heap"
)

func GetHeapSuite(size int) Implementations[algorithms.Heap[int]] {
	less := cmp.Less[int]
	return []Implementation[algorithms.Heap[int]]{
		{
			name: "Gollections/Binary",
			factory: func() algorithms.Heap[int] {
				return heap.NewBinaryHeap(size, less)
			},
		},
		{
			name: "stdlib",
			factory: func() algorithms.Heap[int] {
				return NewStdLibHeap(make([]int, size))
			},
		},
	}
}

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
			b.Run("Library=Gollections/BinaryHeap", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					h := heap.NewBinaryHeap(n+k, less)
					h.Push(initialData...)
					b.StartTimer()
					h.Push(newItems...)
				}
			})

			b.Run("Library=Container/Heap", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					h := NewStdLibHeap(initialData)
					b.StartTimer()
					for _, item := range newItems {
						h.Push(item)
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

	b.Run("Library=Gollections/BinaryHeap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			h := heap.NewBinaryHeapFromSlice(append([]int{}, data...), less)
			b.StartTimer()
			for !h.IsEmpty() {
				h.Pop()
			}
		}
	})

	b.Run("Library=Container/Heap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			h := NewStdLibHeap(append([]int{}, data...))
			b.StartTimer()
			for h.Length() > 0 {
				h.Pop()
			}
		}
	})
}

func BenchmarkBinaryHeap_TopK(b *testing.B) {
	const N = 100_000
	const K = 100

	data := make([]int, N)
	for i := range N {
		data[i] = rand.Int()
	}

	for _, implementation := range GetHeapSuite(K) {
		b.Run("Library="+implementation.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				heap := implementation.factory()
				algorithms.TopK(data, K, heap, func(i1, i2 int) bool { return i1 > i2 })
			}
		})
	}
}
