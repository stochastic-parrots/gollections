package heaps_test

import (
	"fmt"

	"github.com/stochastic-parrots/gollections/pkg/heaps"
)

func ExampleMin() {
	h := heaps.Min[int](5)

	h.Push(10, 50, 5, 1)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleMinFromSlice() {
	data := []int{42, 7, 13, 1, 99}

	h := heaps.MinFromSlice(data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 1
	// Slice: [7 42 13 99]
}

func ExampleMinCloneSlice() {
	data := []int{42, 7, 13, 1, 99}

	h := heaps.MinCloneSlice(data)
	h.Pop()

	fmt.Printf("Length: %d\n", h.Length())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Length: 4
	// Data: [42 7 13 1 99]
}

func ExampleMax() {
	h := heaps.Max[float64](0)
	h.Push(1.5, 10.2, 3.7)

	top, _ := h.Peek()
	fmt.Printf("%.1f\n", top)

	// Output:
	// 10.2
}

func ExampleMaxFromSlice() {
	data := []int{1, 13, 7, 42, 99}

	h := heaps.MaxFromSlice(data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 99
	// Slice: [42 13 7 1]
}

func ExampleMaxCloneSlice() {
	data := []int{1, 13, 7, 42, 99}

	h := heaps.MaxCloneSlice(data)
	h.Pop()

	fmt.Printf("Length: %d\n", h.Length())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Length: 4
	// Data: [1 13 7 42 99]
}

func ExampleNew() {
	byLength := func(a, b string) bool {
		return len(a) < len(b)
	}

	h := heaps.New(0, byLength)
	h.Push("apple", "kiwi", "banana", "pear")

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%s ", val)
	}

	// Output:
	// kiwi pear apple banana
}
