package sortedlist

import (
	"cmp"
	"testing"
)

var skipListBoolSink bool
var skipListIntSink int

func BenchmarkSkipList_AddMany(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)

	b.ReportAllocs()
	for b.Loop() {
		list := NewSkipList(cmp.Compare[int])
		list.Add(data...)
		skipListIntSink = list.Length()
	}
}

func BenchmarkSkipList_AddSingle(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipList(cmp.Compare[int])
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

func BenchmarkSkipList_Find(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		_, skipListBoolSink = list.Find(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkSkipList_Contains(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		skipListBoolSink = list.Contains(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkSkipList_LowerBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		skipListIntSink = list.LowerBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkSkipList_UpperBound(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		skipListIntSink = list.UpperBound(i % arraySortedListBenchmarkSize)
	}
}

func BenchmarkSkipList_Get(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	list := NewSkipListFromSlice(data, cmp.Compare[int])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		value, err := list.Get(i % arraySortedListBenchmarkSize)
		if err != nil {
			b.Fatal(err)
		}
		skipListIntSink = value
	}
}

func BenchmarkSkipList_Remove(b *testing.B) {
	data := descendingInts(arraySortedListBenchmarkSize)
	value := arraySortedListBenchmarkSize / 2

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		list := NewSkipListCloneSlice(data, cmp.Compare[int])
		b.StartTimer()

		skipListBoolSink = list.Remove(value)
	}
}
