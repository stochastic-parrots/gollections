package suites

import (
	"cmp"
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

func GetPriorityMapSuite(size int) datastructs.Implementations[datastructs.PriorityMap[int, float64]] {
	less := cmp.Less[float64]
	return []datastructs.Implementation[datastructs.PriorityMap[int, float64]]{
		{
			Name: "stdlib",
			Factory: func() datastructs.PriorityMap[int, float64] {
				return datastructs.NewStdPriorityMap[int](size, less)
			},
		},
		{
			Name: "Gollections_BinaryHeapPriorityMap",
			Factory: func() datastructs.PriorityMap[int, float64] {
				return prioritymap.NewBinaryPriorityMap[int](size, less)
			},
		},
	}
}

func BenchmarkPriorityMap_Dijkstra(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph(nodes, density)

	for _, implementation := range GetPriorityMapSuite(nodes) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				priorityMap := implementation.Factory()
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
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				priorityMap := implementation.Factory()
				algorithms.Prim(graph, priorityMap)
			}
		})
	}
}
