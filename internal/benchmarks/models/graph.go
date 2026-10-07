package models

import (
	"math/rand/v2"

	"github.com/stochastic-parrots/gollections/internal/shared/constraint"
)

// Edge represents a weighted directed edge in benchmark graph inputs.
type Edge[T constraint.Number] struct {
	To     int
	Weight T
}

// Graph is an adjacency-list graph used by benchmark algorithms.
type Graph[T constraint.Number] [][]Edge[T]

// UndirectedEdge represents an unweighted edge between two graph nodes.
type UndirectedEdge struct {
	From int
	To   int
}

// UndirectedGraph stores nodes and each unweighted undirected edge once.
type UndirectedGraph struct {
	Nodes int
	Edges []UndirectedEdge
}

// WeightedUndirectedEdge represents a weighted edge between two graph nodes.
type WeightedUndirectedEdge[T constraint.Number] struct {
	From   int
	To     int
	Weight T
}

// WeightedUndirectedGraph stores nodes and each weighted undirected edge once.
type WeightedUndirectedGraph[T constraint.Number] struct {
	Nodes int
	Edges []WeightedUndirectedEdge[T]
}

func weight[T constraint.Number](r *rand.Rand) T {
	var zero T

	switch any(zero).(type) {
	case float32, float64:
		return T(r.Float64()*1000 + 1)

	default:
		return T(r.Int64N(1000) + 1)
	}
}

// NewRandomGraph creates a deterministic random weighted directed graph.
func NewRandomGraph[T constraint.Number](nodes int, density float64) Graph[T] {
	pcg := rand.NewPCG(42, 1024)
	r := rand.New(pcg)

	graph := make([][]Edge[T], nodes)
	for i := range nodes {
		for j := range nodes {
			if i != j && r.Float64() < density {
				graph[i] = append(graph[i], Edge[T]{
					To:     j,
					Weight: weight[T](r),
				})
			}
		}
	}
	return graph
}

// NewRandomUndirectedGraph creates a deterministic random unweighted graph.
func NewRandomUndirectedGraph(nodes int, density float64) UndirectedGraph {
	pcg := rand.NewPCG(42, 1024)
	r := rand.New(pcg)
	edges := make([]UndirectedEdge, 0)

	for from := range nodes {
		for to := from + 1; to < nodes; to++ {
			if r.Float64() < density {
				edges = append(edges, UndirectedEdge{From: from, To: to})
			}
		}
	}
	return UndirectedGraph{Nodes: nodes, Edges: edges}
}

// NewRandomWeightedUndirectedGraph creates a deterministic random weighted graph.
func NewRandomWeightedUndirectedGraph[T constraint.Number](nodes int, density float64) WeightedUndirectedGraph[T] {
	pcg := rand.NewPCG(42, 2048)
	r := rand.New(pcg)
	edges := make([]WeightedUndirectedEdge[T], 0)

	for from := range nodes {
		for to := from + 1; to < nodes; to++ {
			if r.Float64() < density {
				edges = append(edges, WeightedUndirectedEdge[T]{
					From:   from,
					To:     to,
					Weight: weight[T](r),
				})
			}
		}
	}
	return WeightedUndirectedGraph[T]{Nodes: nodes, Edges: edges}
}
