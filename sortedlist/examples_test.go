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
	items := sortedlist.NewArray(cmp.Compare[int], 0)
	items.Add(3, 1, 2, 2)

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	first, _ := items.First()
	last, _ := items.Last()
	fmt.Println(first, last)

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	var values []int
	_ = json.Unmarshal([]byte(`[9,7,8]`), &values)
	items = sortedlist.ArrayFrom(cmp.Compare[int], values)
	fmt.Println("Decoded:", slices.Collect(items.All()))

	// Output:
	// [1 2 2 3]
	// [3 2 2 1]
	// 1 3
	// Marshal: [1,2,2,3]
	// Decoded: [7 8 9]
}

func ExampleNewOrderedArray() {
	ascending := sortedlist.OrderedArrayClone(sortedlist.Asc, []int{3, 1, 2, 2})
	descending := sortedlist.OrderedArrayClone(sortedlist.Desc, []int{3, 1, 2, 2})

	fmt.Println(slices.Collect(ascending.All()))
	fmt.Println(slices.Collect(descending.All()))

	// Output:
	// [1 2 2 3]
	// [3 2 2 1]
}

func ExampleArraySortedList_bounds() {
	l := sortedlist.NewOrderedArray[int](sortedlist.Asc, 0)
	l.Add(1, 2, 2, 2, 3)

	start, end := l.EqualRange(2)

	fmt.Println("LowerBound:", l.LowerBound(2))
	fmt.Println("UpperBound:", l.UpperBound(2))
	fmt.Println("EqualRange:", start, end)
	fmt.Println("Count:", l.Count(2))

	// Output:
	// LowerBound: 1
	// UpperBound: 4
	// EqualRange: 1 4
	// Count: 3
}

func ExampleArraySortedList_navigate() {
	l := sortedlist.NewOrderedArray[int](sortedlist.Asc, 0)
	l.Add(10, 20, 30)

	ceiling, ceilingIdx, _ := l.Ceiling(25)
	floor, floorIdx, _ := l.Floor(25)
	higher, higherIdx, _ := l.Higher(20)
	lower, lowerIdx, _ := l.Lower(20)

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
	l := sortedlist.NewOrderedArray[int](sortedlist.Asc, 0)
	l.Add(1, 2, 2, 3, 4)

	fmt.Println(slices.Collect(l.Range(2, 4)))

	// Output:
	// [2 2 3]
}

func ExampleArrayFrom() {
	data := []int{3, 1, 2}
	l := sortedlist.ArrayFrom(cmp.Compare[int], data)

	fmt.Println(slices.Collect(l.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [1 2 3]
}

func ExampleArrayClone() {
	data := []int{3, 1, 2}
	l := sortedlist.ArrayClone(cmp.Compare[int], data)

	fmt.Println(slices.Collect(l.All()))
	fmt.Println(data)

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleArrayFromSeq() {
	source := list.NewArray[int](0)
	source.Append(3, 1, 2)

	l := sortedlist.ArrayFromSeq(cmp.Compare[int], source.All())

	fmt.Println(slices.Collect(l.All()))
	fmt.Println(source.ToSlice())

	// Output:
	// [1 2 3]
	// [3 1 2]
}

func ExampleAsReadonly() {
	mutable := sortedlist.NewArray(cmp.Compare[int], 0)
	mutable.Add(2, 1)

	view := sortedlist.AsReadonly(mutable)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	mutable.Add(0)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	// Output:
	// Readonly view: [1 2]
	// Readonly view: [0 1 2]
}
