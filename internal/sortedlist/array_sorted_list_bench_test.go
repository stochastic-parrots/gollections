package sortedlist

import (
	"cmp"
	"testing"
)

const arraySortedListBenchmarkSize = 100_000

var arraySortedListBoolSink bool
var arraySortedListIntSink int

func BenchmarkArraySortedList_AddMany(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)

	b.ReportAllocs()
	for b.Loop() {
		list := NewArraySortedList(len(data), cmp.Compare[int])
		list.Add(data...)
		arraySortedListIntSink = list.Length()
	}
}

func BenchmarkArraySortedList_AddSingle(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedList(len(data)+1, cmp.Compare[int])
	list.Add(data...)
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
		_, arraySortedListBoolSink = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		arraySortedListBoolSink = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_ContainsMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		arraySortedListBoolSink = list.Contains(-1)
	}
}

func BenchmarkArraySortedList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		arraySortedListIntSink = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		arraySortedListIntSink = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_Count(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		arraySortedListIntSink = list.Count(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkArraySortedList_CountMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		arraySortedListIntSink = list.Count(-1)
	}
}

func BenchmarkArraySortedList_EqualRange(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		start, end := list.EqualRange(i % arraySortedListBenchmarkSize)
		arraySortedListIntSink = end - start
	}
}

func BenchmarkArraySortedList_EqualRangeMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewArraySortedListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		start, end := list.EqualRange(-1)
		arraySortedListIntSink = end - start
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

		arraySortedListBoolSink = list.Remove(value)
	}
}

func descendingInts(size int) []int {
	data := make([]int, size)
	for idx := range size {
		data[idx] = size - idx
	}
	return data
}
