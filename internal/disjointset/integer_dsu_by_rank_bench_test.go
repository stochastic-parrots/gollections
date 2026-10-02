package disjointset

import "testing"

func BenchmarkFlatDisjointSetUnionByRank_FindCompressed(b *testing.B) {
	dsu := NewFlatDisjointSetUnionByRank(0, findForestSize-1)
	buildBalancedFindForest(dsu.Union)
	for component := range findComponents {
		_, _ = dsu.Find(component*findComponentSize + findComponentSize - 1)
	}
	next := 0
	b.ReportAllocs()
	for b.Loop() {
		if next == findComponents {
			next = 0
		}
		_, _ = dsu.Find(next*findComponentSize + findComponentSize - 1)
		next++
	}
}

func BenchmarkFlatDisjointSetUnionByRank_FindPathCompression(b *testing.B) {
	dsu := NewFlatDisjointSetUnionByRank(0, findForestSize-1)
	buildBalancedFindForest(dsu.Union)
	parents := append([]int(nil), dsu.parents...)
	next := 0
	b.ReportAllocs()
	for b.Loop() {
		if next == findComponents {
			b.StopTimer()
			// Restore only the paths compressed by the preceding lookups.
			for base := 0; base < findForestSize; base += findComponentSize {
				for node := base + findComponentSize - 1; node != base; node = parents[node] {
					dsu.parents[node] = parents[node]
				}
			}
			next = 0
			b.StartTimer()
		}
		_, _ = dsu.Find(next*findComponentSize + findComponentSize - 1)
		next++
	}
}

func BenchmarkFlatDisjointSetUnionByRank_RankCompressed(b *testing.B) {
	dsu := NewFlatDisjointSetUnionByRank(0, findForestSize-1)
	buildBalancedFindForest(dsu.Union)
	for component := range findComponents {
		_, _ = dsu.Find(component*findComponentSize + findComponentSize - 1)
	}
	next := 0
	b.ReportAllocs()
	for b.Loop() {
		if next == findComponents {
			next = 0
		}
		_, _ = dsu.Rank(next*findComponentSize + findComponentSize - 1)
		next++
	}
}

func BenchmarkFlatDisjointSetUnionByRank_ConnectedShortPath(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionByRank(0, size-1)
	for idx := 1; idx < size; idx++ {
		dsu.Union(0, idx)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = dsu.Connected(0, size-1)
	}
}

func BenchmarkFlatDisjointSetUnionByRank_UnionSingletons(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionByRank(0, size-1)
	next := 0
	b.ReportAllocs()
	for b.Loop() {
		if next == size {
			b.StopTimer()
			dsu.Reset()
			b.StartTimer()
			next = 0
		}
		dsu.Union(next, next+1)
		next += 2
	}
}

func BenchmarkFlatDisjointSetUnionByRank_UnionExistingSets(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionByRank(0, size-1)
	prepare := func() {
		dsu.Reset()
		for idx := 0; idx < size; idx += 4 {
			dsu.Union(idx, idx+1)
			dsu.Union(idx+2, idx+3)
			dsu.Union(idx+1, idx+3)
		}
	}
	prepare()
	next := 0
	b.ReportAllocs()
	for b.Loop() {
		if next == size {
			b.StopTimer()
			prepare()
			next = 0
			b.StartTimer()
		}
		// Both arguments are nonroot leaves in different four-value sets.
		dsu.Union(next+3, next+7)
		next += 8
	}
}

func BenchmarkFlatDisjointSetUnionByRank_UnionStarBatch(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionByRank(0, size-1)
	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		dsu.Reset()
		b.StartTimer()
		for idx := 1; idx < size; idx++ {
			dsu.Union(0, idx)
		}
	}
}

func BenchmarkFlatDisjointSetUnionByRank_Reset(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionByRank(0, size-1)
	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		for idx := 1; idx < size; idx++ {
			dsu.Union(0, idx)
		}
		b.StartTimer()
		dsu.Reset()
	}
}
