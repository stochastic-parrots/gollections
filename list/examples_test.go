package list_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
)

func ExampleNewArray() {
	list := list.NewArray[int](5)
	list.Append(10, 20, 30)

	val, _ := list.Get(1)
	fmt.Printf("Get(1): %d\n", val)

	_ = list.Set(1, 25)

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(slices.Collect(list.Backward()))

	fmt.Println(list.Contains(50, cmp.Compare[int]))
	list.Append(50)
	fmt.Println(list.Contains(50, cmp.Compare[int]))

	list.Reverse()
	fmt.Println(slices.Collect(list.All()))
	fmt.Println(slices.Collect(list.Backward()))

	_, _ = list.Remove(0)
	_ = list.Insert(1, 89)
	fmt.Println(slices.Collect(list.All()))

	data, _ := json.Marshal(list)
	fmt.Println("Marshal:", string(data))

	input := []byte(`[1,2,3]`)
	_ = json.Unmarshal(input, &list)
	fmt.Println("Unmarshal:", slices.Collect(list.All()))

	// Output:
	// Get(1): 20
	// [10 25 30]
	// [30 25 10]
	// false
	// true
	// [50 30 25 10]
	// [10 25 30 50]
	// [30 89 25 10]
	// Marshal: [30,89,25,10]
	// Unmarshal: [1 2 3]
}

func ExampleNewLinked() {
	list := list.NewLinked[string]()
	list.Append("Go", "is", "fast")

	fmt.Println(slices.Collect(list.All()))
	fmt.Println(slices.Collect(list.Backward()))

	fmt.Println(list.Contains("Go", cmp.Compare[string]))
	fmt.Println(list.Contains("Java", cmp.Compare[string]))

	_ = list.Insert(0, "Java and")
	_, _ = list.Remove(2)
	_ = list.Insert(2, "are")
	fmt.Println(slices.Collect(list.All()))

	data, _ := json.Marshal(list)
	fmt.Println("Marshal:", string(data))

	input := []byte(`["hello", "world"]`)
	_ = json.Unmarshal(input, &list)
	fmt.Println("Unmarshal:", slices.Collect(list.All()))

	// Output:
	// [Go is fast]
	// [fast is Go]
	// true
	// false
	// [Java and Go are fast]
	// Marshal: ["Java and","Go","are","fast"]
	// Unmarshal: [hello world]
}

func ExampleAsReadonly() {
	mutable := list.NewArray[int](0)
	mutable.Append(10, 20)

	data := list.AsReadonly(mutable)

	process := func(view list.Readonly[int]) {
		fmt.Println("Readonly view:", slices.Collect(view.All()))
		// view.Append(30) Compilation Error
	}

	process(data)

	mutable.Append(30)
	process(data)

	// Output:
	// Readonly view: [10 20]
	// Readonly view: [10 20 30]
}
