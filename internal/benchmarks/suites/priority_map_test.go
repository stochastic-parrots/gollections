package suites

import (
	"cmp"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/heap"
)

func GetPriorityMapSuite(size int) Implementations[algorithms.PriorityMap[int, float64]] {
	less := cmp.Less[float64]
	return []Implementation[algorithms.PriorityMap[int, float64]]{
		{
			name: "Gollections/Binary",
			factory: func() algorithms.PriorityMap[int, float64] {
				return heap.NewBinaryPriorityMap[int](size, less)
			},
		},
		{
			name: "stdlib",
			factory: func() algorithms.PriorityMap[int, float64] {
				return NewStdPriorityMap[int](size, less)
			},
		},
	}
}

func BenchmarkPriorityMap_Dijkstra(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph(nodes, density)

	for _, implementation := range GetPriorityMapSuite(nodes) {
		b.Run("Library="+implementation.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				priorityMap := implementation.factory()
				algorithms.Dijkstra(graph, 0, priorityMap)
			}
		})
	}
}

func BenchmarkPriorityMap_Prim(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph(nodes, density)

	for _, implementation := range GetPriorityMapSuite(nodes) {
		b.Run("Library="+implementation.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				priorityMap := implementation.factory()
				algorithms.Prim(graph, priorityMap)
			}
		})
	}
}
