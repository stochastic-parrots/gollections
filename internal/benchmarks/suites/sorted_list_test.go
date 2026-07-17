package suites

import (
	"strconv"
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
				return sortedlist.OrderedArray[int](sortedlist.Asc).New(capacity)
			},
		},
	}
}

func BenchmarkSortedList_Contains(b *testing.B) {
	data := models.NewRandomSlice(10_000_000)
	tests := []struct {
		name   string
		target int
	}{
		{name: "Hit", target: data[len(data)/2]},
		{name: "Miss", target: -1},
	}

	for _, test := range tests {
		b.Run("Case="+test.name, func(b *testing.B) {
			for _, implementation := range getSortedListSuite(len(data)) {
				list := implementation.Factory()
				list.Adds(data...)
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						_ = list.Contains(test.target)
					}
				})
			}
		})
	}
}

func BenchmarkSortedList_RangeQuery(b *testing.B) {
	const size = 1_000_000
	const maxValue = 10_000_000
	const maxWidth = 1_000

	data := models.NewRandomSliceWithMax(size, maxValue)
	from := maxValue / 2
	to := from + maxWidth

	for _, implementation := range getSortedListSuite(len(data)) {
		l := implementation.Factory()
		l.Adds(data...)
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = algorithms.RangeQueryCount(l, from, to)
			}
		})
	}
}

func BenchmarkSortedList_SlidingMedian(b *testing.B) {
	const size = 100_000
	const maxValue = 1_000_000

	data := models.NewRandomSliceWithMax(size, maxValue)
	windowSizes := []int{101, 501, 1_001}

	for _, windowSize := range windowSizes {
		b.Run("WindowSize="+strconv.Itoa(windowSize), func(b *testing.B) {
			windowCount := len(data) - windowSize + 1
			for _, implementation := range getSortedListSuite(windowSize) {
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					list := implementation.Factory()
					b.ReportAllocs()
					for b.Loop() {
						b.StopTimer()
						list.Clear()
						b.StartTimer()
						_ = algorithms.SlidingMedian(data, windowSize, list)
					}
					b.ReportMetric(float64(b.N*windowCount)/b.Elapsed().Seconds(), "windows/s")
				})
			}
		})
	}
}
