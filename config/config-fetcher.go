package config

import (
	"encoding/json"
	"log"
	"os"
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
	Queue     []*any // Fila de Requests TODO: Arrumar o tipo depois
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

// Função que lê o arquivo JSON
func parseServerConfigs(cfg_file string) []ServerUrls {
	fileBytes, err := os.ReadFile(cfg_file)
	if err != nil {
		log.Fatalf("Failed opening jsonfile: %v", err)
	}
	var serverList []ServerUrls
	err = json.Unmarshal(fileBytes, &serverList)
	if err != nil {
		log.Fatalf("Failed parsing jsonfile: %v", err)
	}

	return serverList
}

func MakeServerList(cfg_file string) []*ServerInfo {
	var serverList []*ServerInfo
	urls := parseServerConfigs(cfg_file)
	for i := range len(urls) {
		server := ServerInfo{urls[i], nil, Offline, 0}
		serverList = append(serverList, &server)
	}
	return serverList
}
