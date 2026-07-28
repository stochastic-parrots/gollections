package set_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/set"
)

func ExampleHashSetOf() {
	values := set.HashSetOf[string]().From([]string{"go", "collections", "go"})
	fmt.Println("added:", values.Add("iterators"))
	fmt.Println("added:", values.Add("go"))
	fmt.Println("removed:", values.Removes("collections", "missing", "collections"))

	fmt.Println(slices.Sorted(values.All()))
	fmt.Println(values.Contains("go"))
	fmt.Println(values.Contains("maps"))

	// Output:
	// added: true
	// added: false
	// removed: 1
	// [go iterators]
	// true
	// false
}

func ExampleHashSetBy() {
	type user struct {
		ID    int
		Name  string
		Roles []string
	}

	users := set.HashSetBy(func(value user) int { return value.ID }).New(0)
	added := users.Adds(
		user{ID: 2, Name: "Grace", Roles: []string{"admin"}},
		user{ID: 1, Name: "Ada", Roles: []string{"author"}},
		user{ID: 1, Name: "duplicate"},
	)
	fmt.Println("added:", added)

	ordered := slices.SortedFunc(users.All(), func(left, right user) int {
		return cmp.Compare(left.ID, right.ID)
	})
	for _, value := range ordered {
		fmt.Println(value.ID, value.Name, value.Roles)
	}

	// Output:
	// added: 2
	// 1 Ada [author]
	// 2 Grace [admin]
}

func ExampleHashFactory_Union() {
	factory := set.HashSetOf[int]()
	left := factory.From([]int{1, 2, 3})
	right := factory.From([]int{3, 4})

	union := factory.Union(left, right)

	fmt.Println(slices.Sorted(union.All()))
	fmt.Println(slices.Sorted(left.All()))

	// Output:
	// [1 2 3 4]
	// [1 2 3]
}

func ExampleHashSet_IntersectWith() {
	values := set.HashSetOf[int]().From([]int{1, 2, 3, 4})
	allowed := set.HashSetOf[int]().From([]int{2, 4, 6})

	removed := values.IntersectWith(allowed)

	fmt.Println("removed:", removed)
	fmt.Println(slices.Sorted(values.All()))

	// Output:
	// removed: 2
	// [2 4]
}

func ExampleAsReadonly() {
	mutable := set.HashSetOf[int]().From([]int{1, 2})
	view := set.AsReadonly[int](mutable)

	fmt.Println(slices.Sorted(view.All()))
	mutable.Add(3)
	fmt.Println(slices.Sorted(view.All()))

	// Output:
	// [1 2]
	// [1 2 3]
}
