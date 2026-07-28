package set

import "testing"

type benchmarkRecord struct {
	ID   int
	Data []byte
}

func benchmarkRecordID(value benchmarkRecord) int {
	return value.ID
}

func BenchmarkKeyedHashSet_Adds(b *testing.B) {
	const size = 10_000
	values := make([]benchmarkRecord, size)
	for idx := range values {
		values[idx] = benchmarkRecord{ID: idx}
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewKeyedHashSet(size, benchmarkRecordID)
		b.StartTimer()
		set.Adds(values...)
	}
}

func BenchmarkKeyedHashSet_Contains(b *testing.B) {
	const size = 10_000
	values := make([]benchmarkRecord, size)
	for idx := range values {
		values[idx] = benchmarkRecord{ID: idx}
	}
	set := NewKeyedHashSetFromSlice(values, benchmarkRecordID)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Contains(values[size-1])
	}
}

func BenchmarkKeyedHashSet_Remove(b *testing.B) {
	const size = 10_000
	values := make([]benchmarkRecord, size)
	for idx := range values {
		values[idx] = benchmarkRecord{ID: idx}
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewKeyedHashSetFromSlice(values, benchmarkRecordID)
		b.StartTimer()
		for _, value := range values {
			set.Remove(value)
		}
	}
}

func BenchmarkKeyedHashSet_Removes(b *testing.B) {
	const size = 10_000
	values := make([]benchmarkRecord, size)
	for idx := range values {
		values[idx] = benchmarkRecord{ID: idx}
	}

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewKeyedHashSetFromSlice(values, benchmarkRecordID)
		b.StartTimer()
		set.Removes(values...)
	}
}

func BenchmarkKeyedHashSet_Clone(b *testing.B) {
	const size = 10_000
	values := make([]benchmarkRecord, size)
	for idx := range values {
		values[idx] = benchmarkRecord{ID: idx}
	}
	set := NewKeyedHashSetFromSlice(values, benchmarkRecordID)

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Clone()
	}
}

func BenchmarkKeyedHashSet_Union(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Union(right)
	}
}

func BenchmarkKeyedHashSet_UnionSelf(b *testing.B) {
	set, _ := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Union(set)
	}
}

func BenchmarkKeyedHashSet_UnionWith(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.UnionWith(right)
	}
}

func BenchmarkKeyedHashSet_UnionWithPreallocated(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := NewKeyedHashSet(
			left.Length()+right.Length(),
			benchmarkRecordID,
		)
		set.UnionWith(left)
		b.StartTimer()
		set.UnionWith(right)
	}
}

func BenchmarkKeyedHashSet_Intersection(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Intersection(right)
	}
}

func BenchmarkKeyedHashSet_IntersectionSelf(b *testing.B) {
	set, _ := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Intersection(set)
	}
}

func BenchmarkKeyedHashSet_IntersectWith(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.IntersectWith(right)
	}
}

func BenchmarkKeyedHashSet_Difference(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Difference(right)
	}
}

func BenchmarkKeyedHashSet_DifferenceSelf(b *testing.B) {
	set, _ := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Difference(set)
	}
}

func BenchmarkKeyedHashSet_DifferenceWith(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.DifferenceWith(right)
	}
}

func BenchmarkKeyedHashSet_SymmetricDifference(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.SymmetricDifference(right)
	}
}

func BenchmarkKeyedHashSet_SymmetricDifferenceSelf(b *testing.B) {
	set, _ := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.SymmetricDifference(set)
	}
}

func BenchmarkKeyedHashSet_SymmetricDifferenceWith(b *testing.B) {
	left, right := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		set := left.Clone()
		b.StartTimer()
		set.SymmetricDifferenceWith(right)
	}
}

func BenchmarkKeyedHashSet_Equal(b *testing.B) {
	left, _ := benchmarkKeyedHashSets()
	equal := left.Clone()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.Equal(equal)
	}
}

func BenchmarkKeyedHashSet_EqualSelf(b *testing.B) {
	set, _ := benchmarkKeyedHashSets()

	b.ReportAllocs()
	for b.Loop() {
		_ = set.Equal(set)
	}
}

func BenchmarkKeyedHashSet_IsSubset(b *testing.B) {
	left, superset := benchmarkKeyedHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsSubset(superset)
	}
}

func BenchmarkKeyedHashSet_IsProperSubset(b *testing.B) {
	left, superset := benchmarkKeyedHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsProperSubset(superset)
	}
}

func BenchmarkKeyedHashSet_IsSuperset(b *testing.B) {
	subset, superset := benchmarkKeyedHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = superset.IsSuperset(subset)
	}
}

func BenchmarkKeyedHashSet_IsProperSuperset(b *testing.B) {
	subset, superset := benchmarkKeyedHashSubset()

	b.ReportAllocs()
	for b.Loop() {
		_ = superset.IsProperSuperset(subset)
	}
}

func BenchmarkKeyedHashSet_IsDisjoint(b *testing.B) {
	const size = 10_000
	leftValues := make([]benchmarkRecord, size)
	rightValues := make([]benchmarkRecord, size)
	for idx := range size {
		leftValues[idx] = benchmarkRecord{ID: idx}
		rightValues[idx] = benchmarkRecord{ID: idx + size}
	}
	left := NewKeyedHashSetFromSlice(leftValues, benchmarkRecordID)
	right := NewKeyedHashSetFromSlice(rightValues, benchmarkRecordID)

	b.ReportAllocs()
	for b.Loop() {
		_ = left.IsDisjoint(right)
	}
}

func benchmarkKeyedHashSets() (
	*KeyedHashSet[benchmarkRecord, int],
	*KeyedHashSet[benchmarkRecord, int],
) {
	const size = 10_000
	leftValues := make([]benchmarkRecord, size)
	rightValues := make([]benchmarkRecord, size)
	for idx := range size {
		leftValues[idx] = benchmarkRecord{ID: idx}
		rightValues[idx] = benchmarkRecord{ID: idx + size/2}
	}
	return NewKeyedHashSetFromSlice(leftValues, benchmarkRecordID),
		NewKeyedHashSetFromSlice(rightValues, benchmarkRecordID)
}

func benchmarkKeyedHashSubset() (
	*KeyedHashSet[benchmarkRecord, int],
	*KeyedHashSet[benchmarkRecord, int],
) {
	const size = 10_000
	subsetValues := make([]benchmarkRecord, size)
	supersetValues := make([]benchmarkRecord, size+size/2)
	for idx := range subsetValues {
		subsetValues[idx] = benchmarkRecord{ID: idx}
	}
	for idx := range supersetValues {
		supersetValues[idx] = benchmarkRecord{ID: idx}
	}
	return NewKeyedHashSetFromSlice(subsetValues, benchmarkRecordID),
		NewKeyedHashSetFromSlice(supersetValues, benchmarkRecordID)
}
