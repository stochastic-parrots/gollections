package deque_test

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/deque"
)

func ExampleNewArray() {
	deque := deque.NewArray[int](2)
	deque.Append(2, 3)
	deque.Prepend(0, 1)

	front, _ := deque.Front()
	back, _ := deque.Back()
	fmt.Println(front, back)

	shifted, _ := deque.Shift()
	popped, _ := deque.Pop()
	fmt.Println(shifted, popped)
	fmt.Println(slices.Collect(deque.All()))

	data, _ := json.Marshal(deque)
	fmt.Println("Marshal:", string(data))

	_ = json.Unmarshal([]byte(`[8,9]`), deque)
	fmt.Println("Unmarshal:", slices.Collect(deque.All()))

	// Output:
	// 0 3
	// 0 3
	// [1 2]
	// Marshal: [1,2]
	// Unmarshal: [8 9]
}

func ExampleNewLinked() {
	deque := deque.NewLinked[string]()
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
