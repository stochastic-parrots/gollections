package heap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/shared/ordering"
)

func BenchmarkBinaryHeap_Push(b *testing.B) {
	const capacity = 100_000

	h := NewBinaryHeap(capacity, ordering.Min[int]())
	value := capacity
	b.ReportAllocs()
	for b.Loop() {
		if h.Length() == capacity {
			b.StopTimer()
			h.Clear()
			value = capacity
			b.StartTimer()
		}
		h.Push(value)
		value--
	}
}

func BenchmarkBinaryHeap_Pushes(b *testing.B) {
	const n = 10_000

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

	for _, ratio := range ratios {
		k := int(float64(n) * ratio.val)
		newItems := make([]int, k)

		b.Run(ratio.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				h := NewBinaryHeap(n+k, ordering.Min[int]())
				h.Pushes(initialData...)
				b.StartTimer()
				h.Pushes(newItems...)
			}
		})
	}
}

func BenchmarkBinaryHeap_Pop(b *testing.B) {
	const N = 100_000
	data := make([]int, N)
	for i := range N {
		data[i] = i
	}

	for b.Loop() {
		b.StopTimer()
		h := NewBinaryHeapFromSlice(append([]int{}, data...), ordering.Min[int]())
		b.StartTimer()

		for !h.IsEmpty() {
			h.Pop()
		}
	}

}
