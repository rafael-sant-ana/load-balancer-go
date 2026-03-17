package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	config "github.com/rafael-sant-ana/load-balancer-go/config"
)

// * Sujeito a variações, revisar depois
type healthCheckReponse struct {
	Counter int    `json:"counter"`
	Memory  string `json:"memory_info"`
	CPU     string `json:"cpu_usage"`
}

type GlobalServersInfo struct {
	Total_requests int
	ServerList     []*config.ServerInfo
	BestServer     *config.ServerInfo
	WorstServer    *config.ServerInfo
}

func updateServersStatus(serverList []*config.ServerInfo) {
	var wg sync.WaitGroup

	for _, server := range serverList {
		wg.Add(1)
		go func(s *config.ServerInfo) {
			defer wg.Done()
			serverUrl := server.Info.Url
			healthEndPoint := server.Info.Healthcheck
			url := serverUrl + healthEndPoint

			client := http.Client{}

			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				panic(err)
			}

			resp, err := client.Do(request)
			if err != nil {
				return
			}

			decoder := json.NewDecoder(resp.Body)
			decoder.DisallowUnknownFields()

			var r healthCheckReponse

			e := decoder.Decode(&r)
			if e != nil {
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return
			}
			if r.Counter > 0 {
				server.Status = config.Busy
				return
			}

			server.Status = config.Available

		}(server)
	}
	wg.Wait()
}

func SetupServers() *GlobalServersInfo {
	servers := config.MakeServerList("config/balancer-config.json")
	serversInfo := GlobalServersInfo{0, servers, nil, nil}
	updateServersStatus(servers)
	for _, server := range servers {
		fmt.Printf("%s ServerStatus: %s\n", server.Info.Url, server.Status)
	}
	return &serversInfo
}

func ProcessRequest() *http.Response {
	return nil
}

func main() {
	SetupServers()
}
