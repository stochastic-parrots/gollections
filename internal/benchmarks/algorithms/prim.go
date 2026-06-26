package algorithms

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

// Prim computes the total weight of a minimum spanning tree using the provided priority map.
func Prim[T constraint.Number](graph models.Graph[T], pm datastructs.PriorityMap[int, T]) T {
	var zeroT T
	nodes := len(graph)
	mstWeight := zeroT
	visited := make([]bool, nodes)

	weights := make([]T, nodes)
	inf := infinity[T]()
	for i := range weights {
		weights[i] = inf
	}

	weights[0] = 0
	pm.Set(0, 0)

	for pm.Length() != 0 {
		u, w, _ := pm.Pop()

		if visited[u] {
			continue
		}
		visited[u] = true
		mstWeight += w

		for _, edge := range graph[u] {
			if !visited[edge.To] && edge.Weight < weights[edge.To] {
				weights[edge.To] = edge.Weight
				pm.Set(edge.To, edge.Weight)
			}
		}
	}
	return mstWeight
}
