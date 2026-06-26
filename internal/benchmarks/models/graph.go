package models

import (
	"math/rand/v2"

	"github.com/stochastic-parrots/gollections/constraint"
)

// Edge represents a weighted directed edge in benchmark graph inputs.
type Edge[T constraint.Number] struct {
	To     int
	Weight T
}

// Graph is an adjacency-list graph used by benchmark algorithms.
type Graph[T constraint.Number] [][]Edge[T]

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
