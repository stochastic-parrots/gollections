package sortedlist

import "testing"

var orderedSkipListBoolSink bool
var orderedSkipListIntSink int

func BenchmarkOrderedSkipList_AddMany(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)

	b.ReportAllocs()
	for b.Loop() {
		list := NewOrderedSkipList[int]()
		list.Add(data...)
		orderedSkipListIntSink = list.Length()
	}
}

func BenchmarkOrderedSkipList_AddSingle(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)
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

func BenchmarkOrderedSkipList_Find(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, orderedSkipListBoolSink = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedSkipList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedSkipListBoolSink = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedSkipList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedSkipListIntSink = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedSkipList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		orderedSkipListIntSink = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkOrderedSkipList_Get(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewOrderedSkipListFromSlice(data)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		value, err := list.Get(i % arraySortedListBenchmarkSize)
		if err != nil {
			b.Fatal(err)
		}
		orderedSkipListIntSink = value
	}
}

func BenchmarkOrderedSkipList_Remove(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	value := arraySortedListBenchmarkSize / 2

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		list := NewOrderedSkipListCloneSlice(data)
		b.StartTimer()

		orderedSkipListBoolSink = list.Remove(value)
	}
}
