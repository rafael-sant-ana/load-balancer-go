package types

type ServerHeap []*ServerInfo

func (h ServerHeap) Len() int { return len(h) }
func (h ServerHeap) Less(i, j int) bool {
	if h[i].Status != h[j].Status {
		return h[i].Status < h[j].Status
	}
	return h[i].QueueSize < h[j].QueueSize
}
func (h ServerHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *ServerHeap) Push(x any) {
	*h = append(*h, x.(*ServerInfo))
}

func (h *ServerHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	if x == nil {
		panic("Tried to pop invalid server")
	}
	return x
}

type ServerMaxHeap []*ServerInfo

func (h ServerMaxHeap) Len() int { return len(h) }
func (h ServerMaxHeap) Less(i, j int) bool {
	return h[i].QueueSize > h[j].QueueSize
}
func (h ServerMaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *ServerMaxHeap) Push(x any) {
	*h = append(*h, x.(*ServerInfo))
}
func (h *ServerMaxHeap) Top() *ServerInfo {
	items := *h
	if items.Len() == 0 {
		return nil
	}
	return items[0]
}

func (h *ServerMaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	if x == nil {
		panic("Tried to pop invalid server")
	}
	return x
}
