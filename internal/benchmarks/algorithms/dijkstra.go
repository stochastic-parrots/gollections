package algorithms

import (
	"math"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

func Dijkstra(graph models.Graph, start int, pm PriorityMap[int, float64]) []float64 {
	dist := make([]float64, len(graph))
	for i := range dist {
		dist[i] = math.MaxInt
	}

	dist[start] = 0
	pm.Set(start, 0)

	for pm.Length() > 0 {
		u, d, _ := pm.Pop()
		if d > dist[u] {
			continue
		}

		for _, e := range graph[u] {
			if alt := d + e.Weight; alt < dist[e.To] {
				dist[e.To] = alt
				pm.Set(e.To, alt)
			}
		}
	}
	return dist
}
