package algorithms

import (
	"math"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

func Prim(graph models.Graph, pm PriorityMap[int, float64]) float64 {
	nodes := len(graph)
	mstWeight := 0.0
	visited := make([]bool, nodes)

	weights := make([]float64, nodes)
	for i := range weights {
		weights[i] = math.MaxFloat64
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
