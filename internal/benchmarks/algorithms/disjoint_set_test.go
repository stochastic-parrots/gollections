package algorithms

import (
	"testing"

	"github.com/stochastic-parrots/gollections/disjointset"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stretchr/testify/assert"
)

func TestKruskal(t *testing.T) {
	tests := []struct {
		name             string
		graph            models.WeightedUndirectedGraph[int64]
		wantWeight       int64
		wantDisjointSets int
	}{
		{
			name: "Connected",
			graph: models.WeightedUndirectedGraph[int64]{
				Nodes: 4,
				Edges: []models.WeightedUndirectedEdge[int64]{
					{From: 0, To: 1, Weight: 1},
					{From: 1, To: 2, Weight: 2},
					{From: 0, To: 2, Weight: 4},
					{From: 2, To: 3, Weight: 3},
				},
			},
			wantWeight:       6,
			wantDisjointSets: 1,
		},
		{
			name: "DisconnectedWithIsolatedVertex",
			graph: models.WeightedUndirectedGraph[int64]{
				Nodes: 6,
				Edges: []models.WeightedUndirectedEdge[int64]{
					{From: 0, To: 1, Weight: 1},
					{From: 1, To: 2, Weight: 2},
					{From: 3, To: 4, Weight: 4},
				},
			},
			wantWeight:       7,
			wantDisjointSets: 3,
		},
		{
			name:             "OnlyIsolatedVertices",
			graph:            models.WeightedUndirectedGraph[int64]{Nodes: 3},
			wantWeight:       0,
			wantDisjointSets: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := disjointset.NewIntsRangeBySize(0, test.graph.Nodes-1)

			assert.Equal(t, test.wantWeight, Kruskal(test.graph, set))
			assert.Equal(t, test.wantDisjointSets, set.Disjoints())
		})
	}
}

func TestConnectedComponents(t *testing.T) {
	tests := []struct {
		name           string
		graph          models.UndirectedGraph
		wantComponents int
	}{
		{
			name: "Connected",
			graph: models.UndirectedGraph{
				Nodes: 4,
				Edges: []models.UndirectedEdge{
					{From: 0, To: 1},
					{From: 1, To: 2},
					{From: 0, To: 2},
					{From: 2, To: 3},
				},
			},
			wantComponents: 1,
		},
		{
			name: "DisconnectedWithIsolatedVertex",
			graph: models.UndirectedGraph{
				Nodes: 6,
				Edges: []models.UndirectedEdge{
					{From: 0, To: 1},
					{From: 1, To: 2},
					{From: 3, To: 4},
				},
			},
			wantComponents: 3,
		},
		{
			name:           "OnlyIsolatedVertices",
			graph:          models.UndirectedGraph{Nodes: 3},
			wantComponents: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := disjointset.NewIntsRangeBySize(0, test.graph.Nodes-1)

			assert.Equal(t, test.wantComponents, ConnectedComponents(test.graph, set))
		})
	}
}
