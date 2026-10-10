package ordering

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin_FloatNaN(t *testing.T) {
	nan := math.NaN()
	hasPriority := Min[float64]()

	assert.True(t, hasPriority(nan, 1))
	assert.False(t, hasPriority(1, nan))
	assert.False(t, hasPriority(nan, nan))
}

func TestMax_FloatNaN(t *testing.T) {
	nan := math.NaN()
	hasPriority := Max[float64]()

	assert.True(t, hasPriority(1, nan))
	assert.False(t, hasPriority(nan, 1))
	assert.False(t, hasPriority(nan, nan))
}
