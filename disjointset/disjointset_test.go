package disjointset_test

import (
	"math"
	"slices"
	"testing"

	"github.com/stochastic-parrots/gollections"
	"github.com/stochastic-parrots/gollections/disjointset"
	"github.com/stretchr/testify/assert"
)

type integerID int

func TestNewIntsRangeByRank(t *testing.T) {
	var _ *disjointset.IntsRangeByRank[int] = disjointset.NewIntsRangeByRank(0, 2)
	var _ disjointset.RankDisjointSet[int] = disjointset.NewIntsRangeByRank(0, 2)
	var _ gollections.Collection[int] = disjointset.NewIntsRangeByRank(0, 2)

	set := disjointset.NewIntsRangeByRank(integerID(-1), integerID(1))
	assert.Equal(t, 3, set.Len())
	assert.Equal(t, 3, set.Disjoints())
	assert.False(t, set.IsEmpty())
	assert.Equal(t, []integerID{-1, 0, 1}, slices.Collect(set.All()))
}

func TestNewIntsRangeBySize(t *testing.T) {
	var _ *disjointset.IntsRangeBySize[int] = disjointset.NewIntsRangeBySize(0, 2)
	var _ disjointset.SizeDisjointSet[int] = disjointset.NewIntsRangeBySize(0, 2)
	var _ gollections.Collection[int] = disjointset.NewIntsRangeBySize(0, 2)

	set := disjointset.NewIntsRangeBySize(uint8(254), uint8(255))
	assert.Equal(t, 2, set.Len())
	assert.Equal(t, 2, set.Disjoints())
	assert.False(t, set.IsEmpty())
	assert.Equal(t, []uint8{254, 255}, slices.Collect(set.All()))
}

func TestConstructors_InvalidRange(t *testing.T) {
	for name, construct := range map[string]func(){
		"RankReversed": func() { disjointset.NewIntsRangeByRank(2, 1) },
		"SizeReversed": func() { disjointset.NewIntsRangeBySize(2, 1) },
		"RankTooWide":  func() { disjointset.NewIntsRangeByRank(int64(math.MinInt64), int64(math.MaxInt64)) },
		"SizeTooWide":  func() { disjointset.NewIntsRangeBySize(uint64(0), uint64(math.MaxUint64)) },
	} {
		t.Run(name, func(t *testing.T) { assert.Panics(t, construct) })
	}
}

func TestDisjointSet_PublicContract(t *testing.T) {
	for name, construct := range map[string]func() disjointset.DisjointSet[int]{
		"Rank": func() disjointset.DisjointSet[int] { return disjointset.NewIntsRangeByRank(-1, 2) },
		"Size": func() disjointset.DisjointSet[int] { return disjointset.NewIntsRangeBySize(-1, 2) },
	} {
		t.Run(name, func(t *testing.T) {
			set := construct()
			assert.Equal(t, 4, set.Len())
			assert.False(t, set.IsEmpty())
			assert.Equal(t, []int{-1, 0, 1, 2}, slices.Collect(set.All()))
			var indexes, values []int
			for idx, value := range set.Enumerate() {
				indexes = append(indexes, idx)
				values = append(values, value)
			}
			assert.Equal(t, []int{0, 1, 2, 3}, indexes)
			assert.Equal(t, []int{-1, 0, 1, 2}, values)
			assert.Equal(t, 4, set.Disjoints())
			assert.True(t, set.Union(-1, 0))
			assert.False(t, set.Union(0, -1))
			assert.False(t, set.Union(-2, 1))
			connected, ok := set.Connected(-1, 0)
			assert.True(t, connected)
			assert.True(t, ok)
			connected, ok = set.Connected(-1, 1)
			assert.False(t, connected)
			assert.True(t, ok)
			connected, ok = set.Connected(-2, 0)
			assert.False(t, connected)
			assert.False(t, ok)
			_, ok = set.Find(-2)
			assert.False(t, ok)
			_, ok = set.Find(0)
			assert.True(t, ok)
			assert.Equal(t, 3, set.Disjoints())
			set.Reset()
			assert.Equal(t, 4, set.Disjoints())
			connected, ok = set.Connected(-1, 0)
			assert.False(t, connected)
			assert.True(t, ok)
		})
	}
}

func TestIntsRangeByRank_Rank(t *testing.T) {
	set := disjointset.NewIntsRangeByRank(1, 3)
	rank, ok := set.Rank(0)
	assert.Zero(t, rank)
	assert.False(t, ok)
	assert.True(t, set.Union(1, 2))
	rank, ok = set.Rank(2)
	assert.Equal(t, 1, rank)
	assert.True(t, ok)
}

func TestIntsRangeBySize_Size(t *testing.T) {
	set := disjointset.NewIntsRangeBySize(1, 3)
	size, ok := set.Size(0)
	assert.Zero(t, size)
	assert.False(t, ok)
	assert.True(t, set.Union(1, 2))
	size, ok = set.Size(2)
	assert.Equal(t, 2, size)
	assert.True(t, ok)
}

func TestAsReadonly(t *testing.T) {
	assert.Nil(t, disjointset.AsReadonly[int](nil))

	mutable := disjointset.NewIntsRangeBySize(1, 3)
	view := disjointset.AsReadonly[int](mutable)
	_, mutableView := any(view).(disjointset.DisjointSet[int])
	_, sizeView := any(view).(disjointset.SizeDisjointSet[int])
	assert.False(t, mutableView)
	assert.False(t, sizeView)
	assert.False(t, view.IsEmpty())
	assert.Equal(t, 3, view.Len())
	assert.Equal(t, []int{1, 2, 3}, slices.Collect(view.All()))
	var indexes, values []int
	for idx, value := range view.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}
	assert.Equal(t, []int{0, 1, 2}, indexes)
	assert.Equal(t, []int{1, 2, 3}, values)
	representative, ok := view.Find(1)
	assert.Equal(t, 1, representative)
	assert.True(t, ok)
	assert.Equal(t, 3, view.Disjoints())
	assert.True(t, mutable.Union(1, 2))
	assert.Equal(t, 2, view.Disjoints())
	connected, ok := view.Connected(1, 2)
	assert.True(t, connected)
	assert.True(t, ok)
}
