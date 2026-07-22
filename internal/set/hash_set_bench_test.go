package set

import "testing"

func BenchmarkHashSet_Adds(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSet[int](size)
		b.StartTimer()
		set.Adds(values...)
	}
}

func BenchmarkHashSet_Contains(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}
	set := NewHashSetFromSlice(values)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Contains(size - 1)
	}
}

func BenchmarkHashSet_Remove(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSetFromSlice(values)
		b.StartTimer()
		for _, value := range values {
			set.Remove(value)
		}
	}
}

func BenchmarkHashSet_Removes(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSetFromSlice(values)
		b.StartTimer()
		set.Removes(values...)
	}
}
