package heap

import (
	"testing"
)

func BenchmarkBinaryPriorityMap_Set(b *testing.B) {
	const n = 100_000
	pm := NewBinaryPriorityMap[int](n, MinFunc[int]())
	for i := range n {
		pm.Set(i, i)
	}

	b.ResetTimer()
	for b.Loop() {
		pm.Set(n/2, -1)
	}
}

func BenchmarkBinaryPriorityMap_Pop(b *testing.B) {
	const n = 100_000

	for b.Loop() {
		b.StopTimer()
		pm := NewBinaryPriorityMap[int](n, MinFunc[int]())
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
	pm := NewBinaryPriorityMap[int](n, MinFunc[int]())
	for i := range n {
		pm.Set(i, i)
	}

	b.ResetTimer()
	i := 0
	for b.Loop() {
		key := i % n
		pm.Remove(key)
		pm.Set(key, i)
		i++
	}
}
