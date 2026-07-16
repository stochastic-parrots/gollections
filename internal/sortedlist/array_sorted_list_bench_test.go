package sortedlist

import (
	"cmp"
	"testing"
)

const arraySortedListBenchmarkSize = 100_000

func BenchmarkArraySortedList_Adds(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedList(len(data), cmp.Compare[int])

	b.ReportAllocs()
	for b.Loop() {
		list.Adds(data...)
		_ = list.Length()

		b.StopTimer()
		list.Clear()
		b.StartTimer()
	}
}

func BenchmarkArraySortedList_Add(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedList(len(data)+1, cmp.Compare[int])
	list.Adds(data...)
	value := arraySortedListBenchmarkSize / 2

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		list.Add(value)

		b.StopTimer()
		list.Remove(value)
		b.StartTimer()
	}
}

func BenchmarkArraySortedList_Find(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, _ = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_ContainsMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = list.Contains(-1)
	}
}

func BenchmarkArraySortedList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_Count(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.Count(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_CountMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = list.Count(-1)
	}
}

func BenchmarkArraySortedList_EqualRange(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		start, end := list.EqualRange(i % arraySortedListBenchmarkSize)
		_ = end - start
	}
}

func BenchmarkArraySortedList_EqualRangeMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		start, end := list.EqualRange(-1)
		_ = end - start
	}
}

func BenchmarkArraySortedList_Remove(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	value := arraySortedListBenchmarkSize / 2

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		list := NewArraySortedListCloneSlice(data, cmp.Compare[int])
		b.StartTimer()

		_ = list.Remove(value)
	}
}

func descendingInts(size int) []int {
	data := make([]int, size)
	for idx := range size {
		data[idx] = size - idx
	}
	return data
}
