package algorithms

import "github.com/stochastic-parrots/gollections/internal/benchmarks/models"

// ConnectedComponents returns the number of components after joining every edge.
// The disjoint set must start with one singleton per graph node and is mutated.
func ConnectedComponents(graph models.UndirectedGraph, set DisjointSet) int {
	for _, edge := range graph.Edges {
		set.Union(edge.From, edge.To)
	}
	return set.Disjoints()
}
