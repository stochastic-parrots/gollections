package deque_test

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/deque"
)

func ExampleNewArray() {
	queue := deque.NewArray[int](2)
	queue.Append(2, 3)
	queue.Prepend(0, 1)

	front, _ := queue.Front()
	back, _ := queue.Back()
	fmt.Println(front, back)

	shifted, _ := queue.Shift()
	popped, _ := queue.Pop()
	fmt.Println(shifted, popped)
	fmt.Println(slices.Collect(queue.All()))

	data, _ := json.Marshal(queue)
	fmt.Println("Marshal:", string(data))

	_ = json.Unmarshal([]byte(`[8,9]`), queue)
	fmt.Println("Unmarshal:", slices.Collect(queue.All()))

	// Output:
	// 0 3
	// 0 3
	// [1 2]
	// Marshal: [1,2]
	// Unmarshal: [8 9]
}

func ExampleNewLinked() {
	d := deque.NewLinked[string]()
	d.Append("middle", "back")
	d.Prepend("front")

	fmt.Println(slices.Collect(d.All()))

	front, _ := d.Front()
	back, _ := d.Back()
	fmt.Println(front, back)

	_, _ = d.Shift()
	d.Append("tail")
	fmt.Println(slices.Collect(d.All()))

	// Output:
	// [front middle back]
	// front back
	// [middle back tail]
}

func ExampleArrayFrom() {
	values := []int{10, 20, 30}
	d := deque.ArrayFrom(values)
	d.Prepend(5)
	d.Append(40)

	fmt.Println(d.ToSlice())

	// Output:
	// [5 10 20 30 40]
}

func ExampleArrayClone() {
	values := []int{10, 20, 30}
	d := deque.ArrayClone(values)
	_, _ = d.Shift()

	fmt.Println(d.ToSlice())
	fmt.Println(values)

	// Output:
	// [20 30]
	// [10 20 30]
}

func ExampleLinkedFromSeq() {
	d := deque.LinkedFromSeq(slices.Values([]int{10, 20, 30}))

	fmt.Println(d.ToSlice())

	// Output:
	// [10 20 30]
}

func ExampleAsReadonly() {
	mutable := deque.NewArray[int](0)
	mutable.Append(10, 20)

	view := deque.AsReadonly(mutable)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	mutable.Prepend(5)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	// Output:
	// Readonly view: [10 20]
	// Readonly view: [5 10 20]
}
