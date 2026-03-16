package main

import (
	"encoding/json"
	"log"
	"os"
)

type ServerType struct {
	url         string
	healthcheck string
}

func getServerList() []ServerType {
	fileBytes, err := os.ReadFile("config/balancer-config.json") // revisar o path
	if err != nil {
		log.Fatalf("Failed opening jsonfile: %v", err)
	}
	var serverList []ServerType
	err = json.Unmarshal(fileBytes, &serverList)
	if err != nil {
		log.Fatalf("Failed parsing jsonfile: %v", err)
	}

	return serverList
}
