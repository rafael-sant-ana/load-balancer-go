package types

type CheckServer struct {
	URL     string
	Status  string
	InQueue int
}

type CheckResponse struct {
	RequestsQueueSize int
	Servers           []CheckServer
}
