package models

import "math/rand/v2"

type Edge struct {
	To     int
	Weight float64
}

type Graph [][]Edge

func NewRandomGraph(nodes int, density float64) Graph {
	pcg := rand.NewPCG(42, 1024)
	r := rand.New(pcg)

	graph := make([][]Edge, nodes)
	for i := range nodes {
		for j := range nodes {
			if i != j && r.Float64() < density {
				graph[i] = append(graph[i], Edge{
					To:     j,
					Weight: r.Float64()*1000 + 1,
				})
			}
		}
	}
	return graph
}
