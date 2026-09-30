package disjointset

import "testing"

func BenchmarkFlatDisjointSetUnionBySize_Find(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionBySize(0, size-1)
	for idx := 1; idx < size; idx++ {
		dsu.Union(0, idx)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = dsu.Find(size - 1)
	}
}

func BenchmarkFlatDisjointSetUnionBySize_Connected(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionBySize(0, size-1)
	for idx := 1; idx < size; idx++ {
		dsu.Union(0, idx)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = dsu.Connected(0, size-1)
	}
}

func BenchmarkFlatDisjointSetUnionBySize_Union(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionBySize(0, size-1)
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

func BenchmarkFlatDisjointSetUnionBySize_Reset(b *testing.B) {
	const size = 4096
	dsu := NewFlatDisjointSetUnionBySize(0, size-1)
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
