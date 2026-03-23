package types

import (
	"net/http"
)

// "Enum" que armazena o status do server
type ServerStatus int

const (
	Available ServerStatus = iota
	Busy
	Offline
)

// Armazena o url e o endpoint de healthcheck do server de acordo com o json
type ServerUrls struct {
	Url         string `json:"url"`
	Healthcheck string `json:"healthcheck"`
}

// Teremos uma lista (ou um minheap ordenado pelo tamanho da fila de requests) de Servers com serverinfos
type ServerInfo struct {
	Info      ServerUrls
	Queue     RequestQueue // Fila de Requests TODO: Arrumar o tipo depois
	Status    ServerStatus
	QueueSize int
}

var ServerStatusName = map[ServerStatus]string{
	Available: "Available",
	Busy:      "Busy",
	Offline:   "Offline",
}

func (ss ServerStatus) String() string {
	return ServerStatusName[ss]
}

// * Sujeito a variações, revisar depois
type HealthCheckReponse struct {
	Counter int    `json:"counter"`
	Memory  string `json:"memory_info"`
	CPU     string `json:"cpu_usage"`
}

type GlobalServersInfo struct {
	Total_requests int
	ServerList     []*ServerInfo
	Heap           ServerHeap
	MaxHeap        ServerMaxHeap
}

type RequestEvent struct {
	Request         *http.Request
	ResponseChannel chan *ResponseEvent
}

type ResponseEvent struct {
	Response    *http.Response
	ProcessedBy string
}
