package datastructs

// Implementation names and constructs one benchmark candidate.
type Implementation[T any] struct {
	Name    string
	Factory func() T
}

// Implementations groups benchmark candidates for a suite.
type Implementations[T any] []Implementation[T]
