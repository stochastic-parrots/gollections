package suites

import (
	"cmp"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/heap"
)

func GetHeapSuite(size int, data []int) datastructs.Implementations[datastructs.Heap[int]] {
	less := cmp.Less[int]
	return []datastructs.Implementation[datastructs.Heap[int]]{
		{
			Name: "stdlib",
			Factory: func() datastructs.Heap[int] {
				if data != nil {
					return datastructs.NewStdLibHeapCloneSlice(data)
				}
				return datastructs.NewStdLibHeap(size)
			},
		},
		{
			Name: "Gollections_BinaryHeap",
			Factory: func() datastructs.Heap[int] {
				if data != nil {
					return heap.NewBinaryHeapCloneSlice(data, less)
				}
				return heap.NewBinaryHeap(size, less)
			},
		},
	}
}

func BenchmarkHeaps_Push_Random(b *testing.B) {
	N := 10_000
	data := models.NewRandomSlice(N)
	ratios := []struct {
		name string
		val  float64
	}{
		{"x0.1", 0.1},
		{"x0.5", 0.5},
		{"x1", 1.0},
		{"x10", 10.0},
	}

	for _, ratio := range ratios {
		n := int(float64(N) * ratio.val)
		items := models.NewRandomSliceWithMax(n, N*10)

		b.Run("Size="+ratio.name, func(b *testing.B) {
			for _, implementation := range GetHeapSuite(n, data) {
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					b.ReportAllocs()
					for range b.N {
						b.StopTimer()
						h := implementation.Factory()
						b.StartTimer()
						h.Push(items...)
					}
				})
			}
		})
	}
}

func BenchmarkHeaps_Push_Reverse(b *testing.B) {
	N := 10_000
	data := models.NewReversedSlice(N)
	ratios := []struct {
		name string
		val  float64
	}{
		{"x0.1", 0.1},
		{"x0.5", 0.5},
		{"x1", 1.0},
		{"x10", 10.0},
	}

	for _, ratio := range ratios {
		n := int(float64(N) * ratio.val)
		items := models.NewReversedSliceStartedAt(n, N)

		b.Run("Size="+ratio.name, func(b *testing.B) {
			for _, implementation := range GetHeapSuite(n, data) {
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					b.ReportAllocs()
					for range b.N {
						b.StopTimer()
						h := implementation.Factory()
						b.StartTimer()
						h.Push(items...)
					}
				})
			}
		})
	}
}

func BenchmarkHeaps_Pop(b *testing.B) {
	const N = 100_000
	data := models.NewRandomSlice(N)
	for _, implementation := range GetHeapSuite(N, data) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			h := implementation.Factory()
			for range b.N {
				if h.Length() == 0 {
					b.StopTimer()
					h = implementation.Factory()
					b.StartTimer()
				}
				h.Pop()
			}
		})
	}
}

func BenchmarkHeaps_TopK(b *testing.B) {
	const N = 100_000
	const K = 100
	data := models.NewRandomSlice(N)

	for _, implementation := range GetHeapSuite(K, nil) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				heap := implementation.Factory()
				algorithms.TopK(data, K, heap, func(some, other int) bool { return some > other })
			}
		})
	}
}
