package sortedlist_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
	"github.com/stochastic-parrots/gollections/sortedlist"
)

func ExampleArray() {
	factory := sortedlist.Array(cmp.Compare[int])
	items := factory.New(0)
	items.Adds(3, 1, 2, 2)

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	first, _ := items.First()
	last, _ := items.Last()
	fmt.Println(first, last)

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	var values []int
	_ = json.Unmarshal([]byte(`[9,7,8]`), &values)
	items = factory.From(values)
	fmt.Println("Decoded:", slices.Collect(items.All()))

	// Output:
	// [1 2 2 3]
	// [3 2 2 1]
	// 1 3
	// Marshal: [1,2,2,3]
	// Decoded: [7 8 9]
}

func ExampleOrderedArray() {
	ascending := sortedlist.OrderedArray[int](sortedlist.Ascending).Clone([]int{3, 1, 2, 2})
	descending := sortedlist.OrderedArray[int](sortedlist.Descending).Clone([]int{3, 1, 2, 2})

	fmt.Println(slices.Collect(ascending.All()))
	fmt.Println(slices.Collect(descending.All()))

	// Output:
	// [1 2 2 3]
	// [3 2 2 1]
}

func ExampleArraySortedList_bounds() {
	list := sortedlist.OrderedArray[int](sortedlist.Ascending).New(0)
	list.Adds(1, 2, 2, 2, 3)

	start, end := list.EqualRange(2)

	fmt.Println("LowerBound:", list.LowerBound(2))
	fmt.Println("UpperBound:", list.UpperBound(2))
	fmt.Println("EqualRange:", start, end)
	fmt.Println("Count:", list.Count(2))

	// Output:
	// LowerBound: 1
	// UpperBound: 4
	// EqualRange: 1 4
	// Count: 3
}

func ExampleArraySortedList_navigate() {
	list := sortedlist.OrderedArray[int](sortedlist.Ascending).New(0)
	list.Adds(10, 20, 30)

	ceiling, ceilingIdx, _ := list.Ceiling(25)
	floor, floorIdx, _ := list.Floor(25)
	higher, higherIdx, _ := list.Higher(20)
	lower, lowerIdx, _ := list.Lower(20)

	fmt.Println("Ceiling:", ceiling, ceilingIdx)
	fmt.Println("Floor:", floor, floorIdx)
	fmt.Println("Higher:", higher, higherIdx)
	fmt.Println("Lower:", lower, lowerIdx)

	// Output:
	// Ceiling: 30 2
	// Floor: 20 1
	// Higher: 30 2
	// Lower: 10 0
}

func ExampleArraySortedList_Range() {
	list := sortedlist.OrderedArray[int](sortedlist.Ascending).New(0)
	list.Adds(1, 2, 2, 3, 4)

	fmt.Println(slices.Collect(list.Range(2, 4)))

	// Output:
	// [2 2 3]
}

func ExampleArrayFactory_From() {
	data := []int{3, 1, 2}
	list := sortedlist.Array(cmp.Compare[int]).From(data)

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [1 2 3]
}

func ExampleArrayFactory_Clone() {
	data := []int{3, 1, 2}
	list := sortedlist.Array(cmp.Compare[int]).Clone(data)

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleArrayFactory_FromSeq() {
	source := list.Array[int]().New(0)
	source.Appends(3, 1, 2)

	list := sortedlist.Array(cmp.Compare[int]).FromSeq(source.All())

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(source.ToSlice())

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleAsReadonly() {
	mutable := sortedlist.Array(cmp.Compare[int]).New(0)
	mutable.Adds(2, 1)

	view := sortedlist.AsReadonly(mutable)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	mutable.Add(0)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	// Output:
	// Readonly view: [1 2]
	// Readonly view: [0 1 2]
}
