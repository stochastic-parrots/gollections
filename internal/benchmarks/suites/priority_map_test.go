package suites

import (
	"cmp"
	"testing"

	"github.com/stochastic-parrots/gollections/constraint"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/algorithms"
	ds "github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stochastic-parrots/gollections/internal/benchmarks/models"
	"github.com/stochastic-parrots/gollections/internal/prioritymap"
)

func getPriorityMapSuite[T constraint.Number](size int) ds.Implementations[ds.PriorityMap[int, T]] {
	less := cmp.Less[T]
	return []ds.Implementation[ds.PriorityMap[int, T]]{
		{
			Name: "stdlib",
			Factory: func() ds.PriorityMap[int, T] {
				return ds.NewStdPriorityMap[int](size, less)
			},
		},
		{
			Name: "Gollections_BinaryHeapPriorityMap",
			Factory: func() ds.PriorityMap[int, T] {
				return prioritymap.NewBinaryPriorityMap[int](size, less)
			},
		},
		{
			Name: "Gollections_PairingHeapPriorityMap",
			Factory: func() ds.PriorityMap[int, T] {
				return prioritymap.NewPairingPriorityMapWithCapacity[int](size, less)
			},
		},
	}
}

func getIntPriorityMapSuite[T constraint.Integer](size int) ds.Implementations[ds.PriorityMap[int, T]] {
	suites := getPriorityMapSuite[T](size)
	suites = append(suites, ds.Implementation[ds.PriorityMap[int, T]]{
		Name: "Gollections_RadixHeapPriorityMap",
		Factory: func() ds.PriorityMap[int, T] {
			return prioritymap.NewRadixPriorityMap[int, T](size)
		},
	})
	return suites
}

func BenchmarkPriorityMap_Dijkstra(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph[float64](nodes, density)

	for _, implementation := range getPriorityMapSuite[float64](nodes) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			priorityMap := implementation.Factory()
			for range b.N {
				b.StopTimer()
				priorityMap.Clear()
				b.StartTimer()
				algorithms.Dijkstra(graph, 0, priorityMap)
			}
		})
	}
}

func BenchmarkPriorityMap_DijkstraUint(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph[uint64](nodes, density)

	for _, implementation := range getIntPriorityMapSuite[uint64](nodes) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			priorityMap := implementation.Factory()
			for range b.N {
				b.StopTimer()
				priorityMap.Clear()
				b.StartTimer()
				algorithms.Dijkstra(graph, 0, priorityMap)
			}
		})
	}
}

func BenchmarkPriorityMap_Prim(b *testing.B) {
	const nodes = 1_000
	const density = 0.8
	graph := models.NewRandomGraph[float64](nodes, density)

	for _, implementation := range getPriorityMapSuite[float64](nodes) {
		b.Run("Library="+implementation.Name, func(b *testing.B) {
			b.ReportAllocs()
			priorityMap := implementation.Factory()
			for range b.N {
				b.StopTimer()
				priorityMap.Clear()
				b.StartTimer()
				algorithms.Prim(graph, priorityMap)
			}
		})
	}
}
