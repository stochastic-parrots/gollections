package prioritymap

import "testing"

func BenchmarkRadixPriorityMap_Set(b *testing.B) {
	const n = 100_000
	pm := NewRadixPriorityMap[int, uint64](n)
	for i := range n {
		pm.Set(i, uint64(i))
	}

	b.ResetTimer()
	for b.Loop() {
		pm.Set(n-1, 1)
		b.StopTimer()
		pm.Set(n-1, n-1)
		b.StartTimer()
	}
}

func BenchmarkRadixPriorityMap_Pop(b *testing.B) {
	const n = 100_000

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		pm := NewRadixPriorityMap[int, uint64](n)
		for i := n; i > 0; i-- {
			pm.Set(i, uint64(i))
		}

		b.StartTimer()
		for !pm.IsEmpty() {
			pm.Pop()
		}
	}
}

func BenchmarkRadixPriorityMap_Remove(b *testing.B) {
	const n = 100_000
	pm := NewRadixPriorityMap[int, uint64](n)
	for i := range n {
		pm.Set(i, uint64(i))
	}

	b.ResetTimer()
	i := 0
	for b.Loop() {
		key := (i % (n - 1)) + 1
		pm.Remove(key)
		b.StopTimer()
		pm.Set(key, uint64(i+n))
		i++
		b.StartTimer()
	}
}
