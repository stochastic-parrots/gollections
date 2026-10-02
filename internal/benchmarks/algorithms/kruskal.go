package algorithms

import (
	"cmp"
	"slices"

	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

// DisjointSet provides the operations needed by Kruskal's algorithm.
type DisjointSet interface {
	Union(a, b int) bool
	Disjoints() int
}

// Kruskal returns the weight of a minimum spanning forest and mutates set.
// The disjoint set must start with one singleton per graph node.
func Kruskal[T constraint.Number](graph models.WeightedUndirectedGraph[T], set DisjointSet) T {
	edges := slices.Clone(graph.Edges)
	slices.SortFunc(edges, func(a, b models.WeightedUndirectedEdge[T]) int {
		return cmp.Compare(a.Weight, b.Weight)
	})

	var weight T
	components := set.Disjoints()
	for _, edge := range edges {
		if set.Union(edge.From, edge.To) {
			weight += edge.Weight
			components--
			if components <= 1 {
				break
			}
		}
	}
	return weight
}
