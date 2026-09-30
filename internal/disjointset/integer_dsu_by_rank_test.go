package disjointset

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFlatDisjointSetUnionByRank(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(-2, 1)
	assert.Equal(t, []int{0, 1, 2, 3}, dsu.parents)
	assert.Equal(t, []uint8{0, 0, 0, 0}, dsu.ranks)
	assert.Equal(t, 4, dsu.disjoints)
}

func TestFlatDisjointSetUnionByRank_Find(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(-1, 3)
	representative, ok := dsu.Find(-2)
	assert.Equal(t, 0, representative)
	assert.False(t, ok)
	representative, ok = dsu.Find(4)
	assert.Equal(t, 0, representative)
	assert.False(t, ok)

	dsu.parents[0], dsu.parents[1] = 2, 0
	representative, ok = dsu.Find(0)
	assert.Equal(t, 1, representative)
	assert.True(t, ok)
	assert.Equal(t, 2, dsu.parents[1])
}

func TestFlatDisjointSetUnionByRank_Connected(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(1, 3)
	connected, ok := dsu.Connected(1, 1)
	assert.True(t, connected)
	assert.True(t, ok)
	connected, ok = dsu.Connected(1, 2)
	assert.False(t, connected)
	assert.True(t, ok)
	dsu.Union(1, 2)
	connected, ok = dsu.Connected(1, 2)
	assert.True(t, connected)
	assert.True(t, ok)
	connected, ok = dsu.Connected(0, 1)
	assert.False(t, connected)
	assert.False(t, ok)
	connected, ok = dsu.Connected(1, 4)
	assert.False(t, connected)
	assert.False(t, ok)
}

func TestFlatDisjointSetUnionByRank_Union(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(0, 5)
	assert.False(t, dsu.Union(-1, 0))
	assert.False(t, dsu.Union(0, 6))
	assert.True(t, dsu.Union(0, 1))
	assert.Equal(t, uint8(1), dsu.ranks[0])
	assert.True(t, dsu.Union(2, 0))
	assert.Equal(t, 0, dsu.parents[2])
	assert.True(t, dsu.Union(3, 4))
	assert.True(t, dsu.Union(3, 5))
	assert.True(t, dsu.Union(3, 0))
	assert.Equal(t, 3, dsu.parents[0])
	assert.Equal(t, uint8(2), dsu.ranks[3])
	assert.False(t, dsu.Union(0, 5))
	assert.Equal(t, 1, dsu.Disjoints())
}

func TestFlatDisjointSetUnionByRank_Rank(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(0, 2)
	rank, ok := dsu.Rank(-1)
	assert.Equal(t, 0, rank)
	assert.False(t, ok)
	assert.True(t, dsu.Union(0, 1))
	rank, ok = dsu.Rank(1)
	assert.Equal(t, 1, rank)
	assert.True(t, ok)
}

func TestFlatDisjointSetUnionByRank_Disjoints(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(0, 1)
	assert.Equal(t, 2, dsu.Disjoints())
	dsu.Union(0, 1)
	assert.Equal(t, 1, dsu.Disjoints())
}

func TestFlatDisjointSetUnionByRank_Len(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(4, 6)
	assert.Equal(t, 3, dsu.Len())
	dsu.Union(4, 5)
	assert.Equal(t, 3, dsu.Len())
}

func TestFlatDisjointSetUnionByRank_IsEmpty(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(4, 4)
	assert.False(t, dsu.IsEmpty())
}

func TestFlatDisjointSetUnionByRank_All(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(-1, 1)
	assert.Equal(t, []int{-1, 0, 1}, slices.Collect(dsu.All()))
	var partial []int
	for value := range dsu.All() {
		partial = append(partial, value)
		break
	}
	assert.Equal(t, []int{-1}, partial)
}

func TestFlatDisjointSetUnionByRank_Enumerate(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(-1, 1)
	var indexes, values []int
	for idx, value := range dsu.Enumerate() {
		indexes = append(indexes, idx)
		values = append(values, value)
	}
	assert.Equal(t, []int{0, 1, 2}, indexes)
	assert.Equal(t, []int{-1, 0, 1}, values)
	for idx, value := range dsu.Enumerate() {
		assert.Equal(t, 0, idx)
		assert.Equal(t, -1, value)
		break
	}
}

func TestFlatDisjointSetUnionByRank_Reset(t *testing.T) {
	dsu := NewFlatDisjointSetUnionByRank(0, 3)
	parents, ranks := dsu.parents, dsu.ranks
	dsu.Union(0, 1)
	dsu.Union(2, 3)
	dsu.Union(0, 2)
	dsu.Reset()
	assert.Equal(t, 4, dsu.Disjoints())
	assert.Equal(t, []int{0, 1, 2, 3}, dsu.parents)
	assert.Equal(t, []uint8{0, 0, 0, 0}, dsu.ranks)
	assert.Equal(t, &parents[0], &dsu.parents[0])
	assert.Equal(t, &ranks[0], &dsu.ranks[0])
	assert.True(t, dsu.Union(0, 1))
}
