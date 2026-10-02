package datastructs

// DisjointSet is the union-find capability required by graph algorithms.
type DisjointSet interface {
	Union(a, b int) bool
	Disjoints() int
}
