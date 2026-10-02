package disjointset_test

import (
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/disjointset"
)

func ExampleNewIntsRangeByRank() {
	groups := disjointset.NewIntsRangeByRank(1, 4)
	groups.Union(1, 2)
	connected, ok := groups.Connected(1, 2)
	fmt.Println(connected, ok, groups.Disjoints())
	fmt.Println(slices.Collect(groups.All()))

	// Output:
	// true true 3
	// [1 2 3 4]
}

func ExampleNewIntsRangeBySize() {
	groups := disjointset.NewIntsRangeBySize(1, 4)
	groups.Union(1, 2)
	groups.Union(2, 3)
	size, ok := groups.Size(1)
	fmt.Println(size, ok)
	groups.Reset()
	fmt.Println(groups.Disjoints())

	// Output:
	// 3 true
	// 4
}

func ExampleAsReadonly() {
	groups := disjointset.NewIntsRangeBySize(1, 2)
	view := disjointset.AsReadonly[int](groups)
	fmt.Println(view.Disjoints())
	groups.Union(1, 2)
	fmt.Println(view.Disjoints())

	// Output:
	// 2
	// 1
}
