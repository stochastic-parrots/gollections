package set

import "testing"

type benchmarkRecord struct {
	ID   int
	Data []byte
}

func benchmarkRecordID(value benchmarkRecord) int {
	return value.ID
}

func BenchmarkKeyedHashSet_AddBatch(b *testing.B) {
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
		set.Add(values...)
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

func BenchmarkKeyedHashSet_RemoveBatch(b *testing.B) {
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
		set.Remove(values...)
	}
}
