package main

import (
	"fmt"
	"net/http"

	balancer "github.com/rafael-sant-ana/load-balancer-go/balancer"
	"github.com/rafael-sant-ana/load-balancer-go/types"
)

func main() {
	fmt.Println("Available Servers: ")
	infos := balancer.ServerList

	for _, server := range infos.ServerList {
		if server.Status != types.Offline {
			fmt.Println(server)
		}
	}

	// exemplo de listener de requests
	// for _, server := range ServerList.ServerList {
	// 	channel := make(chan types.RequestEvent)
	// 	server.Channel = &channel
	// 	go func() {
	// 		for event := range *server.Channel {
	// 			fmt.Printf("event: %v\n", event)
	// 		}
	// 	}()
	// fmt.Println(server.Info.Url + " : " + server.Status.String())
	// fmt.Println(server)
	// }

	http.HandleFunc("/", balancer.MakeRequest)
	http.HandleFunc("/check", balancer.CheckServers)
	fmt.Printf("Server is running\n")

	go balancer.ListenQueues()

	http.ListenAndServe(":8090", nil)

}
