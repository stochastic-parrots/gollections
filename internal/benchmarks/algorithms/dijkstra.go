package algorithms

import (
	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

// Dijkstra computes shortest-path distances using the provided priority map.
func Dijkstra[T constraint.Number](graph models.Graph[T], start int, pm datastructs.PriorityMap[int, T]) []T {
	dist := make([]T, len(graph))
	inf := infinity[T]()
	for i := range dist {
		dist[i] = inf
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
