package config

import (
	"encoding/json"
	"log"
	"os"

	Types "github.com/rafael-sant-ana/load-balancer-go/types"
)

// Função que lê o arquivo JSON
func parseServerConfigs(cfg_file string) []Types.ServerUrls {
	fileBytes, err := os.ReadFile(cfg_file)
	if err != nil {
		log.Fatalf("Failed opening jsonfile: %v", err)
	}
	var serverList []Types.ServerUrls
	err = json.Unmarshal(fileBytes, &serverList)
	if err != nil {
		log.Fatalf("Failed parsing jsonfile: %v", err)
	}

	return serverList
}

func MakeServerList(cfg_file string) []*Types.ServerInfo {
	var serverList []*Types.ServerInfo
	urls := parseServerConfigs(cfg_file)
	for _, url := range urls {
		server := Types.ServerInfo{
			Info:      url,
			Queue:     []*Types.RequestEvent{},
			Status:    Types.Offline,
			QueueSize: 0,
		}

		serverList = append(serverList, &server)
	}
	return serverList
}
