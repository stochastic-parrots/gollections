package suites

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/sortedlist"
)

func getSortedListSuite(capacity int) datastructs.Implementations[datastructs.SortedList] {
	return []datastructs.Implementation[datastructs.SortedList]{
		{
			Name: "stdlib",
			Factory: func() datastructs.SortedList {
				return datastructs.NewStdSortedList(capacity)
			},
		},
		{
			Name: "Gollections_ArraySortedList",
			Factory: func() datastructs.SortedList {
				return sortedlist.NewArray(capacity, cmp.Compare[int])
			},
		},
		{
			Name: "Gollections_OrderedArraySortedList",
			Factory: func() datastructs.SortedList {
				return sortedlist.NewOrderedArray[int](capacity)
			},
		},
	}
}

func BenchmarkSortedList_Contains(b *testing.B) {
	data := models.NewRandomSlice(10_000_000)
	slices.Sort(data)
	toSearch := models.NewRandomSliceFrom(data, 100)

	for _, implementation := range getSortedListSuite(len(data)) {
		l := implementation.Factory()
		l.Add(data...)
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				for i := range toSearch {
					if found := l.Contains(toSearch[i]); !found {
						b.FailNow()
					}
				}
			}
		})
	}
}

func BenchmarkSortedList_RangeQueries(b *testing.B) {
	const size = 1_000_000
	const maxValue = 10_000_000
	const queries = 1_000
	const maxWidth = 1_000

	data := models.NewRandomSliceWithMax(size, maxValue)
	rangeQueries := models.NewRandomRangeQueries(queries, maxValue, maxWidth)

	for _, implementation := range getSortedListSuite(len(data)) {
		l := implementation.Factory()
		l.Add(data...)
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				_ = algorithms.RangeQueryCount(l, rangeQueries)
			}
		})
	}
}

func BenchmarkSortedList_SlidingMedian(b *testing.B) {
	const size = 100_000
	const maxValue = 1_000_000
	const windowSize = 1_001

	data := models.NewRandomSliceWithMax(size, maxValue)

	for _, implementation := range getSortedListSuite(windowSize) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			l := implementation.Factory()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				b.StopTimer()
				l.Clear()
				b.StartTimer()
				_ = algorithms.SlidingMedian(data, windowSize, l)
			}
		})
	}
}
