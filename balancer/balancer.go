package balancer

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	config "github.com/rafael-sant-ana/load-balancer-go/config"
	"github.com/rafael-sant-ana/load-balancer-go/types"
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
		Heap:           Types.ServerHeap{},
		MaxHeap:        Types.ServerMaxHeap{},
	}
	updateServersStatus(servers)

	heap.Init(&serversInfo.Heap)
	heap.Init(&serversInfo.MaxHeap)

	for _, server := range servers {
		h := &serversInfo.Heap
		mh := &serversInfo.MaxHeap
		heap.Push(h, server)
		heap.Push(mh, server)
	}
	return &serversInfo
}

var ServerList *Types.GlobalServersInfo

func init() {
	ServerList = SetupServers()
}

func ProcessRequest(s *Types.ServerInfo, r *types.RequestEvent) {
	s.Status = Types.Busy

	// Lógica de processamento aqui simulada por um sleep
	time.Sleep(10 * time.Second)

	//fim
	r.ResponseChannel <- &types.ResponseEvent{Response: nil, ProcessedBy: s.Info.Url}
	s.QueueSize--
	ServerList.Total_requests--

	// Continua processando a fila ou rouba de outro server
	if s.Queue.Top() != nil {
		new_r, err := s.Queue.Dequeue()
		if err != nil {
			panic(err)
		}
		ProcessRequest(s, new_r)
		return
	}
	s.Status = types.Available

	// TODO: Debugar race condition pra habilitar steal task de forma segura
	// if ServerList.MaxHeap.Top().QueueSize != 0 {
	// 	fmt.Println("Roubei!")
	// 	old_s := heap.Pop(&ServerList.MaxHeap).(*types.ServerInfo)
	// 	r, err := old_s.Queue.Dequeue()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ProcessRequest(s, r)
	// }
}

func EnqueueRequest(req *http.Request, list *Types.GlobalServersInfo, rChannel chan *Types.ResponseEvent) {
	server := heap.Pop(&list.Heap).(*Types.ServerInfo)
	if server == nil {
		panic("No server available on the list")
	}

	reqEvent := types.RequestEvent{Request: req, ResponseChannel: rChannel}
	server.Queue.Enqueue(&reqEvent)

	server.QueueSize += 1
	ServerList.Total_requests++

	// TODO: Ajustar a lógica de work steal evitando race conditions
	h := &list.Heap
	// mh := &list.MaxHeap
	// heap.Fix(mh, len(*mh)-1)
	heap.Push(h, server)

}

// Request que retorna como está a lista de request e os servidores
func CheckServers(w http.ResponseWriter, req *http.Request) {
	serversStauts := []Types.CheckServer{}

	for _, server := range ServerList.ServerList {
		info := Types.CheckServer{Status: server.Status.String(), InQueue: server.QueueSize, URL: server.Info.Url}
		serversStauts = append(serversStauts, info)
	}

	r := Types.CheckResponse{
		RequestsQueueSize: ServerList.Total_requests,
		Servers:           serversStauts,
	}

	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	w.Write(resp)
}

func MakeRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request recebido!")
	fmt.Println(req.Method, req.URL.Path)

	// Passar o request pra frente usando o mesmo endpoint porém no servidor disponível.
	// Por enquanto cada request vai simular um request qualquer de 20s nos servers de teste
	ResponseChannel := make(chan *Types.ResponseEvent)
	go EnqueueRequest(req, ServerList, ResponseChannel)

	response := <-ResponseChannel

	// TODO: Fazer o reverse proxy aqui
	fmt.Fprintf(w, "%s", response.ProcessedBy)
}

func ListenQueues() {
	for _, server := range ServerList.ServerList {
		s := server
		go func() {
			for event := range s.Queue.TopChanged {
				if event == nil {
					continue
				}
				if s.Status == types.Offline {
					// TODO: Lançar um 404 e tentar atualizar os servers dnv sempre
					// * Caso tenha ao menos 1 server online nunca chegaremos aqui
					event.ResponseChannel <- &types.ResponseEvent{Response: nil, ProcessedBy: "OFFLINE"}
					continue
				}
				if s.Status != types.Available {
					continue
				}
				r, err := s.Queue.Dequeue()
				if err != nil {
					panic(err)
				}
				ProcessRequest(s, r)
			}
		}()
	}
}
