package main

import (
	"fmt"
	"net/http"

	balancer "github.com/rafael-sant-ana/load-balancer-go/balancer"
)

func main() {
	fmt.Println("Available Servers: ")
	infos := balancer.ServerList

	for _, server := range infos.ServerList {
		fmt.Println(server)
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
	handler := http.HandlerFunc(balancer.MakeRequest)
	fmt.Printf("Server is running\n")
	http.ListenAndServe(":8090", handler)
}
