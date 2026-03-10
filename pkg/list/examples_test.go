package list_test

import (
	"fmt"

	"github.com/stochastic-parrots/gollections/pkg/list"
)

func ExampleNewArray() {
	list := list.NewArray[int](5)
	list.Append(10, 20, 30)

	val, _ := list.Get(1)
	fmt.Printf("Get(1): %d\n", val)

	_ = list.Set(1, 25)

	for v := range list.All() {
		fmt.Printf("%d ", v)
	}

	// Output:
	// Get(1): 20
	// 10 25 30
}

func ExampleNewLinked() {
	list := list.NewLinked[string]()
	list.Append("Go", "is", "fast")

	// O(n)
	list.Reverse()

	for i, v := range list.Enumerate() {
		fmt.Printf("%d:%s ", i, v)
	}

	// Output:
	// 0:fast 1:is 2:Go
}

func ExampleNewDoubleLinked() {
	list := list.NewDoubleLinked[int]()
	list.Append(1, 2, 3)

	// O(1)
	list.Reverse()
	list.Append(4)

	for v := range list.All() {
		fmt.Printf("%d ", v)
	}

	// Output:
	// 3 2 1 4
}
