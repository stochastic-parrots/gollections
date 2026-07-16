package sortedlist

import "testing"

func BenchmarkOrderedArraySortedList_Adds(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedList[int](len(data))

	b.ReportAllocs()
	for b.Loop() {
		list.Adds(data...)
		_ = list.Length()

		b.StopTimer()
		list.Clear()
		b.StartTimer()
	}
}

func BenchmarkOrderedArraySortedList_Add(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedList[int](len(data) + 1)
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

func BenchmarkOrderedArraySortedList_Find(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, _ = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_ContainsMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = list.Contains(-1)
	}
}

func BenchmarkOrderedArraySortedList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_Count(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_ = list.Count(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_CountMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = list.Count(-1)
	}
}

func BenchmarkOrderedArraySortedList_EqualRange(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		start, end := list.EqualRange(i % arraySortedListBenchmarkSize)
		_ = end - start
	}
}

func BenchmarkOrderedArraySortedList_EqualRangeMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		start, end := list.EqualRange(-1)
		_ = end - start
	}
}

func BenchmarkOrderedArraySortedList_Range(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)
	from := arraySortedListBenchmarkSize / 2
	to := from + 100

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		count := 0
		for range list.Range(from, to) {
			count++
		}
		_ = count
	}
}

func BenchmarkOrderedArraySortedList_Remove(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	value := arraySortedListBenchmarkSize / 2

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		list := NewOrderedArraySortedListCloneSlice(data)
		b.StartTimer()

		_ = list.Remove(value)
	}
}
