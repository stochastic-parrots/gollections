package disjointset

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatInterval(t *testing.T) {
	assert.Equal(t, 3, flatInterval(-1, 1))
	assert.Equal(t, 2, flatInterval(uint8(254), uint8(255)))
	assert.PanicsWithValue(t, "disjointset: min must be <= max", func() { flatInterval(1, 0) })
	assert.PanicsWithValue(t, "disjointset: interval is too large", func() {
		flatInterval(int64(math.MinInt64), int64(math.MaxInt64))
	})
}

func TestFlatIndex(t *testing.T) {
	assert.Equal(t, 2, flatIndex(int64(-1), int64(-3)))
	assert.Equal(t, 1, flatIndex(uint64(math.MaxUint64), uint64(math.MaxUint64-1)))
}

func TestFlatValue(t *testing.T) {
	assert.Equal(t, int64(-1), flatValue(2, int64(-3)))
	assert.Equal(t, uint64(math.MaxUint64), flatValue(1, uint64(math.MaxUint64-1)))
}
