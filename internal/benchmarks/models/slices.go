package models

import "math/rand/v2"

// NewRandomSlice creates a slice of random integers.
func NewRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Int()
	}
	return slice
}

// NewRandomSliceWithMax creates a slice of random integers in [0, max).
func NewRandomSliceWithMax(size, maxValue int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.IntN(maxValue)
	}
	return slice
}

// NewReversedSlice creates a descending slice from size to 1.
func NewReversedSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = size - i
	}
	return slice
}

// NewReversedSliceStartedAt creates a descending slice starting at start.
func NewReversedSliceStartedAt(size int, start int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = start - i
	}
	return slice
}
