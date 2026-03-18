package prioritymap

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/comparator"
)

func BenchmarkPairingPriorityMap_Set(b *testing.B) {
	const n = 100_000
	pm := NewPairingPriorityMapWithCapacity[int](n, comparator.Min[int]())
	for i := range n {
		pm.Set(i, i)
	}

	b.ResetTimer()
	for b.Loop() {
		pm.Set(0, n+1)
		b.StopTimer()
		pm.Set(0, 0)
		b.StartTimer()
	}
}

func BenchmarkPairingPriorityMap_Pop(b *testing.B) {
	const n = 100_000

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		pm := NewPairingPriorityMapWithCapacity[int](n, comparator.Min[int]())
		for i := n; i > 0; i-- {
			pm.Set(i, i)
		}

		b.StartTimer()
		for !pm.IsEmpty() {
			pm.Pop()
		}
	}
}

func BenchmarkPairingPriorityMap_Remove(b *testing.B) {
	const n = 100_000
	pm := NewPairingPriorityMapWithCapacity[int](n, comparator.Min[int]())
	for i := range n {
		pm.Set(i, i)
	}

	b.ResetTimer()
	i := 0
	for b.Loop() {
		key := (i % (n - 2)) + 1
		pm.Remove(key)
		b.StopTimer()
		pm.Set(key, i)
		i++
		b.StartTimer()
	}
}
