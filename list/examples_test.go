package list_test

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
)

func ExampleNewArray() {
	items := list.NewArray[int](5)
	items.Appends(10, 20, 30)

	val, _ := items.Get(1)
	fmt.Printf("Get(1): %d\n", val)

	_ = items.Set(1, 25)

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	fmt.Println(items.Contains(func(x int) bool { return x == 50 }))
	items.Append(50)
	fmt.Println(items.Contains(func(x int) bool { return x == 50 }))

	items.Reverse()
	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	_, _ = items.Remove(0)
	_ = items.Insert(1, 89)
	fmt.Println(slices.Collect(items.All()))

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	input := []byte(`[1,2,3]`)
	_ = json.Unmarshal(input, items)
	fmt.Println("Unmarshal:", slices.Collect(items.All()))

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
	items := list.NewLinked[string]()
	items.Appends("Go", "is", "fast")

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	fmt.Println(items.Contains(func(x string) bool { return x == "Go" }))
	fmt.Println(items.Contains(func(x string) bool { return x == "Java" }))

	_ = items.Insert(0, "Java and")
	_, _ = items.Remove(2)
	_ = items.Insert(2, "are")
	fmt.Println(slices.Collect(items.All()))

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	input := []byte(`["hello", "world"]`)
	_ = json.Unmarshal(input, items)
	fmt.Println("Unmarshal:", slices.Collect(items.All()))

	// Output:
	// [Go is fast]
	// [fast is Go]
	// true
	// false
	// [Java and Go are fast]
	// Marshal: ["Java and","Go","are","fast"]
	// Unmarshal: [hello world]
}

func ExampleArrayFrom() {
	values := []int{10, 20, 30}
	list := list.ArrayFrom[int](values)

	list.Append(40)
	fmt.Println(list.ToSlice())

	// Output:
	// [10 20 30 40]
}

func ExampleArrayClone() {
	values := []int{10, 20, 30}
	list := list.ArrayClone[int](values)
	_ = list.Set(0, 100)

	fmt.Println(list.ToSlice())
	fmt.Println(values)

	// Output:
	// [100 20 30]
	// [10 20 30]
}

func ExampleLinkedFromSeq() {
	list := list.LinkedFromSeq[int](slices.Values([]int{10, 20, 30}))

	fmt.Println(list.ToSlice())

	// Output:
	// [10 20 30]
}

func ExampleAsReadonly() {
	mutable := list.NewArray[int](0)
	mutable.Appends(10, 20)

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

func ExampleList_optionalCapabilities() {
	var items list.List[int] = list.ArrayFrom([]int{1, 2})

	if stringer, ok := items.(fmt.Stringer); ok {
		fmt.Println(stringer.String())
	}
	if marshaler, ok := items.(json.Marshaler); ok {
		data, err := marshaler.MarshalJSON()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
	if unmarshaler, ok := items.(json.Unmarshaler); ok {
		if err := unmarshaler.UnmarshalJSON([]byte(`[3,4]`)); err != nil {
			panic(err)
		}
	}
	fmt.Println(items.ToSlice())

	// Output:
	// [1 2]
	// [1,2]
	// [3 4]
}
