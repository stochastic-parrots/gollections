package algorithms

func TopK(data []int, k int, heap Heap[int], hasPriority func(int, int) bool) {
	for _, x := range data {
		if heap.Length() < k {
			heap.Push(x)
		} else {
			root, _ := heap.Peek()
			if hasPriority(x, root) {
				heap.Pop()
				heap.Push(x)
			}
		}
	}
}
