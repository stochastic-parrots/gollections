package heap_test

import (
	"cmp"
	"fmt"

	"github.com/stochastic-parrots/gollections/heap"
)

func ExampleNewBinary() {
	byLength := func(a, b string) bool {
		return len(a) < len(b)
	}

	h := heap.NewBinary(byLength, 0)
	h.Push("apple", "kiwi", "banana", "pear")

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%s ", val)
	}

	// Output:
	// kiwi pear apple banana
}

func ExampleBinaryFrom() {
	data := []int{10, 50, 5, 1}
	h := heap.BinaryFrom(cmp.Less[int], data)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleBinaryClone() {
	data := []int{10, 50, 5, 1}
	h := heap.BinaryClone(cmp.Less[int], data)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleNewOrderedBinary_min() {
	h := heap.NewOrderedBinary[int](heap.Min, 5)

	h.Push(10, 50, 5, 1)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleBinaryFrom_min() {
	data := []int{42, 7, 13, 1, 99}

	h := heap.OrderedBinaryFrom[int](heap.Min, data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 1
	// Slice: [7 42 13 99]
}

func ExampleBinaryClone_min() {
	data := []int{42, 7, 13, 1, 99}

	h := heap.OrderedBinaryClone[int](heap.Min, data)
	h.Pop()

	fmt.Printf("Len: %d\n", h.Len())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Len: 4
	// Data: [42 7 13 1 99]
}

func ExampleNewOrderedBinary_max() {
	h := heap.NewOrderedBinary[float64](heap.Max, 0)
	h.Push(1.5)
	h.Push(10.2, 3.7)

	top, _ := h.Peek()
	fmt.Printf("%.1f\n", top)

	// Output:
	// 10.2
}

func ExampleBinaryFrom_max() {
	data := []int{1, 13, 7, 42, 99}

	h := heap.OrderedBinaryFrom[int](heap.Max, data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 99
	// Slice: [42 13 7 1]
}

func ExampleBinaryClone_max() {
	data := []int{1, 13, 7, 42, 99}

	h := heap.OrderedBinaryClone[int](heap.Max, data)
	h.Pop()

	fmt.Printf("Len: %d\n", h.Len())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Len: 4
	// Data: [1 13 7 42 99]
}
