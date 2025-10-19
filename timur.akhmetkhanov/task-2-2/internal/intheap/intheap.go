package intheap

type IntHeap []int //nolint:recvcheck

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= h.Len() || j >= h.Len() {
		panic("index out of range intheap")
	}

	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	if i < 0 || j < 0 || i >= h.Len() || j >= h.Len() {
		panic("index out of range in intheap")
	}

	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("Invalid type")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
