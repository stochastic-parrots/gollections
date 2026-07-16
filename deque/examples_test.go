package deque_test

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/deque"
)

func ExampleArray() {
	queue := deque.Array[int]().New(2)
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

func ExampleLinked() {
	deque := deque.Linked[string]().New()
	deque.Append("middle", "back")
	deque.Prepend("front")

	fmt.Println(slices.Collect(deque.All()))

	front, _ := deque.Front()
	back, _ := deque.Back()
	fmt.Println(front, back)

	_, _ = deque.Shift()
	deque.Append("tail")
	fmt.Println(slices.Collect(deque.All()))

	// Output:
	// [front middle back]
	// front back
	// [middle back tail]
}

func ExampleArrayFactory_From() {
	values := []int{10, 20, 30}
	deque := deque.Array[int]().From(values)
	deque.Prepend(5)
	deque.Append(40)

	fmt.Println(deque.ToSlice())

	// Output:
	// [5 10 20 30 40]
}

func ExampleArrayFactory_Clone() {
	values := []int{10, 20, 30}
	deque := deque.Array[int]().Clone(values)
	_, _ = deque.Shift()

	fmt.Println(deque.ToSlice())
	fmt.Println(values)

	// Output:
	// [20 30]
	// [10 20 30]
}

func ExampleLinkedFactory_FromSeq() {
	deque := deque.Linked[int]().FromSeq(slices.Values([]int{10, 20, 30}))

	fmt.Println(deque.ToSlice())

	// Output:
	// [10 20 30]
}

func ExampleAsReadonly() {
	mutable := deque.Array[int]().New(0)
	mutable.Append(10, 20)

	view := deque.AsReadonly(mutable)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	mutable.Prepend(5)
	fmt.Println("Readonly view:", slices.Collect(view.All()))

	// Output:
	// Readonly view: [10 20]
	// Readonly view: [5 10 20]
}
