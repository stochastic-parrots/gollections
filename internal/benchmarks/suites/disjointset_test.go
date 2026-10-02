package suites

import (
	"testing"

	"github.com/stochastic-parrots/gollections/disjointset"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
)

type benchmarkDisjointSet interface {
	algorithms.DisjointSet
	Reset()
}

type disjointSetImplementation struct {
	Name    string
	Factory func() benchmarkDisjointSet
}

func getDisjointSetSuite(nodes int) []disjointSetImplementation {
	return []disjointSetImplementation{
		{
			Name: "Gollections_IntsRangeByRank",
			Factory: func() benchmarkDisjointSet {
				return disjointset.NewIntsRangeByRank(0, nodes-1)
			},
		},
		{
			Name: "Gollections_IntsRangeBySize",
			Factory: func() benchmarkDisjointSet {
				return disjointset.NewIntsRangeBySize(0, nodes-1)
			},
		},
	}
}

func BenchmarkDisjointSet_Kruskal(b *testing.B) {
	const nodes = 1024
	for _, graphCase := range []struct {
		name    string
		density float64
	}{
		{name: "Sparse", density: 0.0015},
		{name: "Dense", density: 0.08},
	} {
		graph := models.NewRandomWeightedUndirectedGraph[int64](nodes, graphCase.density)
		b.Run(graphCase.name, func(b *testing.B) {
			for _, implementation := range getDisjointSetSuite(nodes) {
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					set := implementation.Factory()
					b.ReportAllocs()
					for b.Loop() {
						b.StopTimer()
						set.Reset()
						b.StartTimer()
						_ = algorithms.Kruskal(graph, set)
					}
				})
			}
		})
	}
}

func BenchmarkDisjointSet_ConnectedComponents(b *testing.B) {
	const nodes = 1024
	for _, graphCase := range []struct {
		name    string
		density float64
	}{
		{name: "Sparse", density: 0.0015},
		{name: "Dense", density: 0.08},
	} {
		graph := models.NewRandomUndirectedGraph(nodes, graphCase.density)
		b.Run(graphCase.name, func(b *testing.B) {
			for _, implementation := range getDisjointSetSuite(nodes) {
				b.Run("Library="+implementation.Name, func(b *testing.B) {
					set := implementation.Factory()
					b.ReportAllocs()
					for b.Loop() {
						b.StopTimer()
						set.Reset()
						b.StartTimer()
						_ = algorithms.ConnectedComponents(graph, set)
					}
				})
			}
		})
	}
}
