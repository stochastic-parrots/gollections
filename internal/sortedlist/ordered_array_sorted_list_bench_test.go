package sortedlist

import "testing"

var orderedArraySortedListBoolSink bool
var orderedArraySortedListIntSink int

func BenchmarkOrderedArraySortedList_AddMany(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)

	b.ReportAllocs()
	for b.Loop() {
		list := NewOrderedArraySortedList[int](len(data))
		list.Add(data...)
		orderedArraySortedListIntSink = list.Length()
	}
}

func BenchmarkOrderedArraySortedList_AddSingle(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedList[int](len(data) + 1)
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

func BenchmarkOrderedArraySortedList_Find(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, orderedArraySortedListBoolSink = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedArraySortedListBoolSink = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_ContainsMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		orderedArraySortedListBoolSink = list.Contains(-1)
	}
}

func BenchmarkOrderedArraySortedList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedArraySortedListIntSink = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedArraySortedListIntSink = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_Count(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedArraySortedListIntSink = list.Count(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedArraySortedList_CountMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		orderedArraySortedListIntSink = list.Count(-1)
	}
}

func BenchmarkOrderedArraySortedList_EqualRange(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		start, end := list.EqualRange(i % arraySortedListBenchmarkSize)
		orderedArraySortedListIntSink = end - start
	}
}

func BenchmarkOrderedArraySortedList_EqualRangeMissing(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedArraySortedListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		start, end := list.EqualRange(-1)
		orderedArraySortedListIntSink = end - start
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
		orderedArraySortedListIntSink = count
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

		orderedArraySortedListBoolSink = list.Remove(value)
	}
}
