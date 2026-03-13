package heap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/comparator"
)

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
			for b.Loop() {
				b.StopTimer()
				h := NewBinaryHeap(n+k, comparator.Min[int]())
				h.Push(initialData...)
				b.StartTimer()
				h.Push(newItems...)
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
		h := NewBinaryHeapFromSlice(append([]int{}, data...), comparator.Min[int]())
		b.StartTimer()

		for !h.IsEmpty() {
			h.Pop()
		}
	}

}
