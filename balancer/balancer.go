package balancer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

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
		Heap:           []*Types.ServerInfo{},
	}
	updateServersStatus(servers)

	for _, server := range servers {
		serversInfo.Heap.Push(server)
	}
	return &serversInfo
}

var ServerList *Types.GlobalServersInfo

func init() {
	ServerList = SetupServers()
}

func SendRequest(req *http.Request, list *Types.GlobalServersInfo, ResponseChannel chan Types.ResponseEvent) {
	// Aqui ao inves de pegar o "bestServer" vou pegar o root do minHeap, atualizar ele e devolver pro mesmo
	server := list.Heap.Pop()
	if server == nil {
		panic("No server available on the list")
	}
	fmt.Println("Received request to " + server.Info.Url)
	fmt.Println(&ResponseChannel)

	// TODO: Criar uma funcao pra processar o request e retornar no fim.
	// ? Talvez cada server com um listener da respectiva fila? Mandar sempre o request junto com um channel pra response!

	server.QueueSize++
	list.Heap.Push(server)

	time.Sleep(3 * time.Second)
	ResponseChannel <- Types.ResponseEvent{ProcessedBy: server.Info.Url}

	server.QueueSize--
}

func MakeRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request recebido!")
	fmt.Println(req.Method, req.URL.Path)
	// Passar o request pra frente usando o mesmo endpoint porém no servidor disponível.
	// Por enquanto cada request vai simular um request qualquer de 20s nos servers de teste
	ResponseChannel := make(chan Types.ResponseEvent)
	go SendRequest(req, ServerList, ResponseChannel)

	response := <-ResponseChannel
	fmt.Print("Resposta: ")
	fmt.Println(response)
	// Escreve a response
	fmt.Fprintf(w, "Request processed by %s", response.ProcessedBy)
}
