package heap_test

import (
	"cmp"
	"fmt"

	"github.com/stochastic-parrots/gollections/heap"
)

func ExampleBinary() {
	byLength := func(a, b string) bool {
		return len(a) < len(b)
	}

	h := heap.Binary(byLength).New(0)
	h.Push("apple", "kiwi", "banana", "pear")

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%s ", val)
	}

	// Output:
	// kiwi pear apple banana
}

func ExampleBinaryFactory_From() {
	data := []int{10, 50, 5, 1}
	h := heap.Binary(cmp.Less[int]).From(data)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleBinaryFactory_Clone() {
	data := []int{10, 50, 5, 1}
	h := heap.Binary(cmp.Less[int]).Clone(data)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleOrderedBinary_min() {
	h := heap.OrderedBinary[int](heap.Min).New(5)

	h.Push(10, 50, 5, 1)

	for !h.IsEmpty() {
		val, _ := h.Pop()
		fmt.Printf("%d ", val)
	}

	// Output:
	// 1 5 10 50
}

func ExampleBinaryFactory_From_min() {
	data := []int{42, 7, 13, 1, 99}

	h := heap.OrderedBinary[int](heap.Min).From(data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 1
	// Slice: [7 42 13 99]
}

func ExampleBinaryFactory_Clone_min() {
	data := []int{42, 7, 13, 1, 99}

	h := heap.OrderedBinary[int](heap.Min).Clone(data)
	h.Pop()

	fmt.Printf("Length: %d\n", h.Length())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Length: 4
	// Data: [42 7 13 1 99]
}

func ExampleOrderedBinary_max() {
	h := heap.OrderedBinary[float64](heap.Max).New(0)
	h.Push(1.5, 10.2, 3.7)

	top, _ := h.Peek()
	fmt.Printf("%.1f\n", top)

	// Output:
	// 10.2
}

func ExampleBinaryFactory_From_max() {
	data := []int{1, 13, 7, 42, 99}

	h := heap.OrderedBinary[int](heap.Max).From(data)

	val, _ := h.Pop()
	fmt.Printf("Pop: %d\n", val)
	fmt.Printf("Slice: %v\n", data[:4])

	// Output:
	// Pop: 99
	// Slice: [42 13 7 1]
}

func ExampleBinaryFactory_Clone_max() {
	data := []int{1, 13, 7, 42, 99}

	h := heap.OrderedBinary[int](heap.Max).Clone(data)
	h.Pop()

	fmt.Printf("Length: %d\n", h.Length())
	fmt.Printf("Data: %v\n", data)

	// Output:
	// Length: 4
	// Data: [1 13 7 42 99]
}
