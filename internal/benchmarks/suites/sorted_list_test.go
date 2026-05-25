package suites

import (
	"cmp"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/list"
	"github.com/stochastic-parrots/gollections/internal/sortedlist"
)

const sortedListBenchmarkSize = 100_000

var sortedListBoolSink bool
var sortedListIntSink int

func BenchmarkSortedList_Build_Random(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := sortedlist.NewArraySortedList(sortedListBenchmarkSize, cmp.Compare[int])
			list.Add(data...)
		}
	})

	b.Run("Library=Gollections_ArrayListThenSort", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewArrayList[int](sortedListBenchmarkSize)
			list.Append(data...)
			_ = sortedlist.NewArraySortedListFromSlice(list.ToSlice(), cmp.Compare[int])
		}
	})

	b.Run("Library=Gollections_DoubleLinkedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewDoubleLinkedList[int]()
			list.Append(data...)
		}
	})
}

func BenchmarkSortedList_AddOneByOne_Random(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := sortedlist.NewArraySortedList[int](sortedListBenchmarkSize, cmp.Compare[int])
			for _, value := range data {
				list.Add(value)
			}
		}
	})

	b.Run("Library=Gollections_ArrayList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewArrayList[int](sortedListBenchmarkSize)
			for _, value := range data {
				list.Append(value)
			}
		}
	})

	b.Run("Library=Gollections_DoubleLinkedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewDoubleLinkedList[int]()
			for _, value := range data {
				list.Append(value)
			}
		}
	})
}

func BenchmarkSortedList_Contains(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)
	targets := sortedListLookupTargets(data)

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		list := sortedlist.NewArraySortedListFromSlice(append([]int(nil), data...), cmp.Compare[int])

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortedListBoolSink = list.Contains(targets[i%len(targets)])
		}
	})

	b.Run("Library=Gollections_ArrayList", func(b *testing.B) {
		b.ReportAllocs()
		list := list.NewArrayList[int](sortedListBenchmarkSize)
		list.Append(data...)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortedListBoolSink = list.Contains(targets[i%len(targets)], cmp.Compare[int])
		}
	})

	b.Run("Library=Gollections_DoubleLinkedList", func(b *testing.B) {
		b.ReportAllocs()
		list := list.NewDoubleLinkedList[int]()
		list.Append(data...)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortedListBoolSink = list.Contains(targets[i%len(targets)], cmp.Compare[int])
		}
	})
}

func BenchmarkSortedList_Remove_Middle(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)
	target := data[sortedListBenchmarkSize/2]

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := sortedlist.NewArraySortedListFromSlice(append([]int(nil), data...), cmp.Compare[int])
			sortedListBoolSink = list.Remove(target)
		}
	})

	b.Run("Library=Gollections_ArrayList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewArrayList[int](sortedListBenchmarkSize)
			list.Append(data...)

			idx, ok := list.Find(target, cmp.Compare[int])
			if ok {
				_, _ = list.Remove(idx)
			}
			sortedListBoolSink = ok
		}
	})

	b.Run("Library=Gollections_DoubleLinkedList", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			list := list.NewDoubleLinkedList[int]()
			list.Append(data...)

			idx, ok := list.Find(target, cmp.Compare[int])
			if ok {
				_, _ = list.Remove(idx)
			}
			sortedListBoolSink = ok
		}
	})
}

func BenchmarkSortedList_All(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		list := sortedlist.NewArraySortedListFromSlice(append([]int(nil), data...), cmp.Compare[int])

		b.ResetTimer()
		for range b.N {
			total := 0
			for value := range list.All() {
				total += value
			}
			sortedListIntSink = total
		}
	})

	b.Run("Library=Gollections_ArrayList", func(b *testing.B) {
		b.ReportAllocs()
		list := list.NewArrayList[int](sortedListBenchmarkSize)
		list.Append(data...)

		b.ResetTimer()
		for range b.N {
			total := 0
			for value := range list.All() {
				total += value
			}
			sortedListIntSink = total
		}
	})

	b.Run("Library=Gollections_DoubleLinkedList", func(b *testing.B) {
		b.ReportAllocs()
		list := list.NewDoubleLinkedList[int]()
		list.Append(data...)

		b.ResetTimer()
		for range b.N {
			total := 0
			for value := range list.All() {
				total += value
			}
			sortedListIntSink = total
		}
	})
}

func BenchmarkSortedList_BuildOnceLookupMany(b *testing.B) {
	data := models.NewRandomSlice(sortedListBenchmarkSize)
	targets := sortedListLookupTargets(data)

	b.Run("Library=Gollections_ArraySortedList", func(b *testing.B) {
		b.ReportAllocs()
		list := sortedlist.NewArraySortedListFromSlice(append([]int(nil), data...), cmp.Compare[int])

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortedListBoolSink = list.Contains(targets[i%len(targets)])
		}
	})

	b.Run("Library=Gollections_ArrayList", func(b *testing.B) {
		b.ReportAllocs()
		list := list.NewArrayList[int](sortedListBenchmarkSize)
		list.Append(data...)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sortedListBoolSink = list.Contains(targets[i%len(targets)], cmp.Compare[int])
		}
	})
}

func sortedListLookupTargets(data []int) []int {
	targets := make([]int, 128)
	for i := range targets {
		if i%2 == 0 {
			targets[i] = data[(i*7919)%len(data)]
			continue
		}
		targets[i] = -i - 1
	}
	return targets
}
