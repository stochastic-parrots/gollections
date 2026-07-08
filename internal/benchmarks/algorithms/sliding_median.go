package algorithms

import "github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"

// SlidingMedian sums the median of each fixed-size sliding window.
func SlidingMedian(data []int, windowSize int, list datastructs.SortedList) int {
	sum := 0
	medianIdx := windowSize / 2

	for idx, value := range data {
		list.Add(value)
		if idx+1 < windowSize {
			continue
		}

		median, err := list.Get(medianIdx)
		if err != nil {
			panic(err)
		}
		sum += median
		list.Remove(data[idx-windowSize+1])
	}

	return sum
}
