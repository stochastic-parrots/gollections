package models

import "math/rand/v2"

func NewRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Int()
	}
	return slice
}

func NewRandomSliceWithMax(size, max int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.IntN(max)
	}
	return slice
}

func NewReversedSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = size - i
	}
	return slice
}

func NewReversedSliceStartedAt(size int, start int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = start - i
	}
	return slice
}
