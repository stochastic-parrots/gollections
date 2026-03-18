package prioritymap_test

import (
	"cmp"
	"fmt"

	"github.com/stochastic-parrots/gollections/prioritymap"
)

func ExampleNewBinary() {
	pm := prioritymap.NewBinary[string](3, cmp.Less[int])
	pm.Set("some", 10)
	pm.Set("other", 2)
	pm.Set("some", 1)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - some: age - 1
	// user - other: age - 2
}

func ExampleMinBinary() {
	pm := prioritymap.MinBinary[string, int](3)
	pm.Set("some", 3)
	pm.Set("other", 2)

	if _, exists := pm.Get("another"); !exists {
		fmt.Print("user - another: not exists\n")
	}

	pm.Set("another", 10)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - another: not exists
	// user - other: age - 2
	// user - some: age - 3
	// user - another: age - 10
}

func ExampleMaxBinary() {
	pm := prioritymap.MaxBinary[string, int](3)
	pm.Set("some", 10)
	pm.Set("other", 2)

	if _, exists := pm.Get("another"); !exists {
		fmt.Print("user - another: not exists\n")
	}

	pm.Set("another", 12)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - another: not exists
	// user - another: age - 12
	// user - some: age - 10
	// user - other: age - 2
}

func ExampleNewPairing() {
	pm := prioritymap.NewPairing[string](3, cmp.Less[int])
	pm.Set("some", 10)
	pm.Set("other", 2)
	pm.Set("some", 1)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - some: age - 1
	// user - other: age - 2
}

func ExampleMinPairing() {
	pm := prioritymap.MinPairing[string, int](3)
	pm.Set("some", 3)
	pm.Set("other", 2)

	if _, exists := pm.Get("another"); !exists {
		fmt.Print("user - another: not exists\n")
	}

	pm.Set("another", 10)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - another: not exists
	// user - other: age - 2
	// user - some: age - 3
	// user - another: age - 10
}

func ExampleMaxPairing() {
	pm := prioritymap.MaxPairing[string, int](3)
	pm.Set("some", 10)
	pm.Set("other", 2)

	if _, exists := pm.Get("another"); !exists {
		fmt.Print("user - another: not exists\n")
	}

	pm.Set("another", 12)

	for !pm.IsEmpty() {
		key, value, _ := pm.Pop()
		fmt.Printf("user - %s: age - %d\n", key, value)
	}

	// Output:
	// user - another: not exists
	// user - another: age - 12
	// user - some: age - 10
	// user - other: age - 2
}

func ExampleAsReadonly() {
	pm := prioritymap.MinPairing[string, int](10)
	pm.Set("Critical Bug", 1)
	pm.Set("Feature Request", 10)
	pm.Set("Documentation Update", 5)

	view := prioritymap.AsReadonly(pm)

	// view.Pop() => compiler(MissingFieldOrMethod)

	fmt.Println(view.Length())
	fmt.Println(view.IsEmpty())

	priority, _ := view.Get("Feature Request")
	fmt.Println(priority)

	task, priority, _ := view.Peek()
	fmt.Println(task, priority)

	// Output:
	// 3
	// false
	// 10
	// Critical Bug 1
}
