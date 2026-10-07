package algorithms

import (
	"cmp"
	"slices"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/shared/constraint"
)

// Kruskal returns the weight of a minimum spanning forest and mutates set.
// The disjoint set must start with one singleton per graph node.
func Kruskal[T constraint.Number](graph models.WeightedUndirectedGraph[T], set datastructs.DisjointSet) T {
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
