package algorithms

import "github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"

func TopK(data []int, k int, heap datastructs.Heap[int], hasPriority func(int, int) bool) {
	for _, x := range data {
		if heap.Length() < k {
			heap.Push(x)
			continue
		}

		if root, ok := heap.Peek(); ok && hasPriority(x, root) {
			heap.Replace(x)
		}
	}
}
