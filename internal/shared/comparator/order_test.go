package comparator

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin_FloatNaN(t *testing.T) {
	nan := math.NaN()
	min := Min[float64]()

	assert.True(t, min(nan, 1))
	assert.False(t, min(1, nan))
	assert.False(t, min(nan, nan))
}

func TestMax_FloatNaN(t *testing.T) {
	nan := math.NaN()
	max := Max[float64]()

	assert.True(t, max(1, nan))
	assert.False(t, max(nan, 1))
	assert.False(t, max(nan, nan))
}
