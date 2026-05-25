package sortedlist_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stochastic-parrots/gollections/sortedlist"
)

func ExampleNewArray() {
	list := sortedlist.NewArray(0, cmp.Compare[int])
	list.Add(3, 1, 2, 2)

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(slices.Collect(list.Backward()))

	first, _ := list.First()
	last, _ := list.Last()
	fmt.Println(first, last)

	data, _ := json.Marshal(list)
	fmt.Println("Marshal:", string(data))

	_ = json.Unmarshal([]byte(`[9,7,8]`), list)
	fmt.Println("Unmarshal:", slices.Collect(list.All()))

	// Output:
	// [1 2 2 3]
	// [3 2 2 1]
	// 1 3
	// Marshal: [1,2,2,3]
	// Unmarshal: [7 8 9]
}

func ExampleArrayFrom() {
	data := []int{3, 1, 2}
	list := sortedlist.ArrayFrom(data, cmp.Compare[int])

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [1 2 3]
}

func ExampleArrayClone() {
	data := []int{3, 1, 2}
	list := sortedlist.ArrayClone(data, cmp.Compare[int])

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleArrayFromSeq() {
	source := list.NewArray[int](0)
	source.Append(3, 1, 2)

	list := sortedlist.ArrayFromSeq(source.All(), cmp.Compare[int])

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(source.ToSlice())

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleAsReadonly() {
	mutable := sortedlist.NewArray(0, cmp.Compare[int])
	mutable.Add(2, 1)

	view := sortedlist.AsReadonly(mutable)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	mutable.Add(0)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	// Output:
	// Readonly view: [1 2]
	// Readonly view: [0 1 2]
}
