package set

import "testing"

func BenchmarkHashSet_Adds(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSet[int](size)
		b.StartTimer()
		set.Adds(values...)
	}
}

func BenchmarkHashSet_Contains(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}
	set := NewHashSetFromSlice(values)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Contains(size - 1)
	}
}

func BenchmarkHashSet_Remove(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSetFromSlice(values)
		b.StartTimer()
		for _, value := range values {
			set.Remove(value)
		}
	}
}

func BenchmarkHashSet_Removes(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSetFromSlice(values)
		b.StartTimer()
		set.Removes(values...)
	}
}

func BenchmarkHashSet_Clone(b *testing.B) {
	const size = 10_000
	values := make([]int, size)
	for idx := range values {
		values[idx] = idx
	}
	set := NewHashSetFromSlice(values)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Clone()
	}
}

func BenchmarkHashSet_Union(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Union(right)
	}
}

func BenchmarkHashSet_UnionSelf(b *testing.B) {
	set, _ := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Union(set)
	}
}

func BenchmarkHashSet_UnionWith(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.UnionWith(right)
	}
}

func BenchmarkHashSet_UnionWithPreallocated(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewHashSet[int](left.Length() + right.Length())
		set.UnionWith(left)
		b.StartTimer()
		set.UnionWith(right)
	}
}

func BenchmarkHashSet_Intersection(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Intersection(right)
	}
}

func BenchmarkHashSet_IntersectionSelf(b *testing.B) {
	set, _ := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Intersection(set)
	}
}

func BenchmarkHashSet_IntersectWith(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.IntersectWith(right)
	}
}

func BenchmarkHashSet_Difference(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Difference(right)
	}
}

func BenchmarkHashSet_DifferenceSelf(b *testing.B) {
	set, _ := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Difference(set)
	}
}

func BenchmarkHashSet_DifferenceWith(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.DifferenceWith(right)
	}
}

func BenchmarkHashSet_SymmetricDifference(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.SymmetricDifference(right)
	}
}

func BenchmarkHashSet_SymmetricDifferenceSelf(b *testing.B) {
	set, _ := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.SymmetricDifference(set)
	}
}

func BenchmarkHashSet_SymmetricDifferenceWith(b *testing.B) {
	left, right := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.SymmetricDifferenceWith(right)
	}
}

func BenchmarkHashSet_Equal(b *testing.B) {
	left, _ := benchmarkHashSets()
	equal := left.Clone()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Equal(equal)
	}
}

func BenchmarkHashSet_EqualSelf(b *testing.B) {
	set, _ := benchmarkHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Equal(set)
	}
}

func BenchmarkHashSet_IsSubset(b *testing.B) {
	left, superset := benchmarkHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsSubset(superset)
	}
}

func BenchmarkHashSet_IsProperSubset(b *testing.B) {
	left, superset := benchmarkHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsProperSubset(superset)
	}
}

func BenchmarkHashSet_IsSuperset(b *testing.B) {
	subset, superset := benchmarkHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = superset.IsSuperset(subset)
	}
}

func BenchmarkHashSet_IsProperSuperset(b *testing.B) {
	subset, superset := benchmarkHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = superset.IsProperSuperset(subset)
	}
}

func BenchmarkHashSet_IsDisjoint(b *testing.B) {
	const size = 10_000
	leftValues := make([]int, size)
	rightValues := make([]int, size)
	for idx := range size {
		leftValues[idx] = idx
		rightValues[idx] = idx + size
	}
	left := NewHashSetFromSlice(leftValues)
	right := NewHashSetFromSlice(rightValues)

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsDisjoint(right)
	}
}

func benchmarkHashSets() (*HashSet[int], *HashSet[int]) {
	const size = 10_000
	leftValues := make([]int, size)
	rightValues := make([]int, size)
	for idx := range size {
		leftValues[idx] = idx
		rightValues[idx] = idx + size/2
	}
	return NewHashSetFromSlice(leftValues), NewHashSetFromSlice(rightValues)
}

func benchmarkHashSubset() (*HashSet[int], *HashSet[int]) {
	const size = 10_000
	subsetValues := make([]int, size)
	supersetValues := make([]int, size+size/2)
	for idx := range subsetValues {
		subsetValues[idx] = idx
	}
	for idx := range supersetValues {
		supersetValues[idx] = idx
	}
	return NewHashSetFromSlice(subsetValues), NewHashSetFromSlice(supersetValues)
}
