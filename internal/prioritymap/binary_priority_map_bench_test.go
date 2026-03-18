package prioritymap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/comparator"
)

func BenchmarkBinaryPriorityMap_Set(b *testing.B) {
	const n = 100_000
	pm := NewBinaryPriorityMap[int](n, comparator.Min[int]())
	for i := range n {
		pm.Set(i, i)
	}

	b.ResetTimer()
	for b.Loop() {
		pm.Set(n-1, -1)

		b.StopTimer()
		pm.Set(n-1, n-1)
		b.StartTimer()
	}
}

func BenchmarkBinaryPriorityMap_Pop(b *testing.B) {
	const n = 100_000

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		pm := NewBinaryPriorityMap[int](n, comparator.Min[int]())
		for i := range n {
			pm.Set(i, i)
		}

		b.StartTimer()
		for !pm.IsEmpty() {
			pm.Pop()
		}
	}
}

func BenchmarkBinaryPriorityMap_Remove(b *testing.B) {
	const n = 100_000
	pm := NewBinaryPriorityMap[int](n, comparator.Min[int]())
	for i := range n {
		pm.Set(i, i)
	}

	key, _, _ := pm.Peek()
	b.ResetTimer()
	for b.Loop() {
		pm.Remove(key)

		b.StopTimer()
		pm.Set(key, n+1)
		key, _, _ = pm.Peek()
		b.StartTimer()
	}
}
