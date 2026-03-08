package skip

import (
	"math"
)

const (
	// defaultMaxLevel = 18 is a balanced "sweet spot".
	// It supports up to ~262,000 elements with O(log n) efficiency.
	defaultMaxLevel = 18
	minLevel        = 16 // Minimum height to ensure O(log n) for datasets up to 65k elements.
	maxLevel        = 32 // Maximum height to support up to 4 Billion elements.
)

// Node represents a single element in the Skip List.
type Node[T any] struct {
	value T
	next  []*Node[T]
}

// NewNode initializes a node with a specific height.
func NewNode[T any](x T, levels int) *Node[T] {
	return &Node[T]{x, make([]*Node[T], levels)}
}

// SkipList is a probabilistic data structure that allows O(log n) search,
// insertion, and deletion by maintaining multiple layers of linked lists.
type SkipList[T any] struct {
	first        *Node[T]
	hasPriority  func(a, b T) bool // Comparison function (e.g., a < b)
	currentLevel int               // Highest level currently in use
	maxLevel     int               // Absolute maximum height allowed
}

func NewSkipList[T any](hasPriority func(T, T) bool) *SkipList[T] {
	var zero T
	first := NewNode(zero, defaultMaxLevel)
	return &SkipList[T]{first, hasPriority, 0, defaultMaxLevel}
}

func NewSkipListWithExpectedSize[T any](expectedSize int, hasPriority func(T, T) bool) *SkipList[T] {
	var zero T
	maxLevelByExpectedSize := int(math.Ceil(math.Log2(float64(expectedSize))))
	calculatedLevel := min(max(maxLevelByExpectedSize, minLevel), maxLevel)

	first := NewNode(zero, calculatedLevel)
	return &SkipList[T]{first, hasPriority, 0, calculatedLevel}
}

/*func (l *SkipList[T]) append(x T) {
	r := rand.Uint32()
	levels := bits.TrailingZeros32(r)

	if levels > l.maxLevel {
		levels = l.maxLevel - 1
	}

	NewNode(x, levels)
}*/
