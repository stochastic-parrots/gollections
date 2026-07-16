package algorithms

import (
	"testing"

	"github.com/stochastic-parrots/gollections/internal/benchmarks/datastructs"
	"github.com/stretchr/testify/assert"
)

func TestRangeQueryCount(t *testing.T) {
	list := datastructs.NewStdSortedList(5)
	list.Adds(1, 2, 2, 3, 5)

	assert.Equal(t, 3, RangeQueryCount(list, 2, 5))
}
