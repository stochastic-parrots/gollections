package disjointset

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFlatDisjointSetUnionBySize(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(uint8(254), uint8(255))
	assert.Equal(t, []int{-1, -1}, dsu.parents)
	assert.Equal(t, 2, dsu.disjoints)
}

func TestFlatDisjointSetUnionBySize_Find(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(-1, 3)
	representative, ok := dsu.Find(-2)
	assert.Equal(t, 0, representative)
	assert.False(t, ok)
	representative, ok = dsu.Find(4)
	assert.Equal(t, 0, representative)
	assert.False(t, ok)

	dsu.parents[0], dsu.parents[1], dsu.parents[2] = 2, 0, -3
	representative, ok = dsu.Find(0)
	assert.Equal(t, 1, representative)
	assert.True(t, ok)
	assert.Equal(t, 2, dsu.parents[1])
}

func TestFlatDisjointSetUnionBySize_Connected(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(1, 3)
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

func TestFlatDisjointSetUnionBySize_Union(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(0, 4)
	assert.False(t, dsu.Union(-1, 0))
	assert.False(t, dsu.Union(0, 5))
	assert.True(t, dsu.Union(0, 1))
	assert.Equal(t, -2, dsu.parents[0])
	assert.True(t, dsu.Union(2, 0))
	assert.Equal(t, 0, dsu.parents[2])
	assert.Equal(t, -3, dsu.parents[0])
	assert.True(t, dsu.Union(3, 4))
	assert.True(t, dsu.Union(3, 0))
	assert.Equal(t, 0, dsu.parents[3])
	assert.Equal(t, -5, dsu.parents[0])
	assert.False(t, dsu.Union(4, 2))
	assert.Equal(t, 1, dsu.Disjoints())
}

func TestFlatDisjointSetUnionBySize_Size(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(0, 2)
	size, ok := dsu.Size(-1)
	assert.Equal(t, 0, size)
	assert.False(t, ok)
	size, ok = dsu.Size(2)
	assert.Equal(t, 1, size)
	assert.True(t, ok)
	assert.True(t, dsu.Union(0, 1))
	size, ok = dsu.Size(1)
	assert.Equal(t, 2, size)
	assert.True(t, ok)
}

func TestFlatDisjointSetUnionBySize_Disjoints(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(0, 1)
	assert.Equal(t, 2, dsu.Disjoints())
	dsu.Union(0, 1)
	assert.Equal(t, 1, dsu.Disjoints())
}

func TestFlatDisjointSetUnionBySize_Len(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(4, 6)
	assert.Equal(t, 3, dsu.Len())
	dsu.Union(4, 5)
	assert.Equal(t, 3, dsu.Len())
}

func TestFlatDisjointSetUnionBySize_IsEmpty(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(4, 4)
	assert.False(t, dsu.IsEmpty())
}

func TestFlatDisjointSetUnionBySize_All(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(-1, 1)
	assert.Equal(t, []int{-1, 0, 1}, slices.Collect(dsu.All()))
	var partial []int
	for value := range dsu.All() {
		partial = append(partial, value)
		break
	}
	assert.Equal(t, []int{-1}, partial)
}

func TestFlatDisjointSetUnionBySize_Enumerate(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(-1, 1)
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

func TestFlatDisjointSetUnionBySize_Reset(t *testing.T) {
	dsu := NewFlatDisjointSetUnionBySize(0, 3)
	parents := dsu.parents
	dsu.Union(0, 1)
	dsu.Union(2, 3)
	dsu.Union(0, 2)
	dsu.Reset()
	assert.Equal(t, 4, dsu.Disjoints())
	assert.Equal(t, []int{-1, -1, -1, -1}, dsu.parents)
	assert.Equal(t, &parents[0], &dsu.parents[0])
	assert.True(t, dsu.Union(0, 1))
}
