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

func (h *ServerHeap) Push(x *ServerInfo) {
	*h = append(*h, x)
}

func (h *ServerHeap) Pop() *ServerInfo {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	if x == nil {
		panic("Tried to pop invalid server")
	}
	return x
}
