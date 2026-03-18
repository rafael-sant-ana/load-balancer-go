package balancer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	config "github.com/rafael-sant-ana/load-balancer-go/config"
	Types "github.com/rafael-sant-ana/load-balancer-go/types"
)

func updateServersStatus(serverList []*Types.ServerInfo) {
	var wg sync.WaitGroup

	for _, server := range serverList {
		wg.Add(1)
		go func(s *Types.ServerInfo) {
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

			var r Types.HealthCheckReponse

			e := decoder.Decode(&r)
			if e != nil {
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return
			}
			if r.Counter > 0 {
				server.Status = Types.Busy
				return
			}

			server.Status = Types.Available

		}(server)
	}
	wg.Wait()
}

func SetupServers() *Types.GlobalServersInfo {
	servers := config.MakeServerList("config/balancer-config.json")
	if servers == nil {
		panic("Erro ao criar a lista de servidores!")
	}
	serversInfo := Types.GlobalServersInfo{
		Total_requests: 0,
		ServerList:     servers,
		BestServer:     servers[0],
		WorstServer:    servers[0],
	}

	updateServersStatus(servers)
	for _, server := range servers {
		fmt.Printf("%s ServerStatus: %s\n", server.Info.Url, server.Status)
	}
	return &serversInfo
}

func EnqueueRequest() {}

func ProcessRequest(originalRequest *http.Request, list *Types.GlobalServersInfo) *http.Response {
	returnChannel := make(chan Types.RequestEvent)

	result := <-returnChannel
	fmt.Println(result.ProcessedBy)
	return result.Response
}

func main() {
	SetupServers()
}
