package list_test

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/list"
)

func ExampleArray() {
	factory := list.Array[int]()
	items := factory.New(5)
	items.Append(10, 20, 30)

	val, _ := items.Get(1)
	fmt.Printf("Get(1): %d\n", val)

	_ = items.Set(1, 25)

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	fmt.Println(items.Contains(50, cmp.Compare[int]))
	items.Append(50)
	fmt.Println(items.Contains(50, cmp.Compare[int]))

	items.Reverse()
	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	_, _ = items.Remove(0)
	_ = items.Insert(1, 89)
	fmt.Println(slices.Collect(items.All()))

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	input := []byte(`[1,2,3]`)
	var values []int
	_ = json.Unmarshal(input, &values)
	items = factory.From(values)
	fmt.Println("Decoded:", slices.Collect(items.All()))

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
	// Decoded: [1 2 3]
}

func ExampleLinked() {
	factory := list.Linked[string]()
	items := factory.New()
	items.Append("Go", "is", "fast")

	fmt.Println(slices.Collect(items.All()))
	fmt.Println(slices.Collect(items.Backward()))

	fmt.Println(items.Contains("Go", cmp.Compare[string]))
	fmt.Println(items.Contains("Java", cmp.Compare[string]))

	_ = items.Insert(0, "Java and")
	_, _ = items.Remove(2)
	_ = items.Insert(2, "are")
	fmt.Println(slices.Collect(items.All()))

	data, _ := json.Marshal(items)
	fmt.Println("Marshal:", string(data))

	input := []byte(`["hello", "world"]`)
	var values []string
	_ = json.Unmarshal(input, &values)
	items = factory.From(values)
	fmt.Println("Decoded:", slices.Collect(items.All()))

	// Output:
	// [Go is fast]
	// [fast is Go]
	// true
	// false
	// [Java and Go are fast]
	// Marshal: ["Java and","Go","are","fast"]
	// Decoded: [hello world]
}

func ExampleArrayFactory_From() {
	values := []int{10, 20, 30}
	list := list.Array[int]().From(values)

	list.Append(40)
	fmt.Println(list.ToSlice())

	// Output:
	// [10 20 30 40]
}

func ExampleArrayFactory_Clone() {
	values := []int{10, 20, 30}
	list := list.Array[int]().Clone(values)
	_ = list.Set(0, 100)

	fmt.Println(list.ToSlice())
	fmt.Println(values)

	// Output:
	// [100 20 30]
	// [10 20 30]
}

func ExampleLinkedFactory_FromSeq() {
	list := list.Linked[int]().FromSeq(slices.Values([]int{10, 20, 30}))

	fmt.Println(list.ToSlice())

	// Output:
	// [10 20 30]
}

func ExampleAsReadonly() {
	mutable := list.Array[int]().New(0)
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
