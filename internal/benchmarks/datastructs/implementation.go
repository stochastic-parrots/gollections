package datastructs

type Implementation[T any] struct {
	Name    string
	Factory func() T
}

type Implementations[T any] []Implementation[T]
