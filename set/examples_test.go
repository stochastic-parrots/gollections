package set_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/stochastic-parrots/gollections/set"
)

func ExampleNewHashSet() {
	values := set.HashSetFrom[string]([]string{"go", "collections", "go"})
	fmt.Println("added:", values.Add("iterators"))
	fmt.Println("added:", values.Add("go"))
	fmt.Println("removed:", values.Remove("collections", "missing", "collections"))

	fmt.Println(slices.Sorted(values.All()))
	fmt.Println(values.Contains("go"))
	fmt.Println(values.Contains("maps"))

	// Output:
	// added: 1
	// added: 0
	// removed: 1
	// [go iterators]
	// true
	// false
}

func ExampleNewKeyedHashSet() {
	type user struct {
		ID    int
		Name  string
		Roles []string
	}

	users := set.NewKeyedHashSet(func(value user) int { return value.ID }, 0)
	added := users.Add(
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

func ExampleEqualBy() {
	type user struct {
		ID   int
		Name string
	}

	byName := func(value user) string { return value.Name }
	left := set.KeyedHashSetFrom(byName, []user{{ID: 1, Name: "Ada"}})
	right := set.KeyedHashSetFrom(byName, []user{{ID: 1, Name: "Ada Lovelace"}})

	fmt.Println(set.EqualBy(func(value user) int { return value.ID }, left, right))

	// Output:
	// true
}

func ExampleHashSet_Union() {
	left := set.HashSetFrom([]int{1, 2, 3})
	right := set.HashSetFrom([]int{3, 4})

	union := left.Union(right)

	fmt.Println(slices.Sorted(union.All()))
	fmt.Println(slices.Sorted(left.All()))

	// Output:
	// [1 2 3 4]
	// [1 2 3]
}

func ExampleSourceFromSlice() {
	result := set.NewHashSet[int](0).Union(set.SourceFromSlice([]int{1, 2, 2}))

	fmt.Println(slices.Sorted(result.All()))

	// Output:
	// [1 2]
}

func ExampleHashSet_IntersectWith() {
	values := set.HashSetFrom([]int{1, 2, 3, 4})
	allowed := set.HashSetFrom([]int{2, 4, 6})

	removed := values.IntersectWith(allowed)

	fmt.Println("removed:", removed)
	fmt.Println(slices.Sorted(values.All()))

	// Output:
	// removed: 2
	// [2 4]
}

func ExampleAsReadonly() {
	mutable := set.HashSetFrom[int]([]int{1, 2})
	view := set.AsReadonly[int](mutable)

	fmt.Println(slices.Sorted(view.All()))
	mutable.Add(3)
	fmt.Println(slices.Sorted(view.All()))

	// Output:
	// [1 2]
	// [1 2 3]
}
