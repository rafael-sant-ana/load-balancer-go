package balancer

import (
	"fmt"
	"sync"
)

type Container struct {
	mu           sync.Mutex
	bestServer   string
	currentValue int
}

func (c *Container) updateBestServer() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Lógica de comparação aqui
}

func get_best_server(urls []string) {
	var endpoint string = "k"

	// var bestserver = None
	for i := range len(urls) {
		fmt.Println(urls[i])
		fmt.Println(endpoint)
		// lock
		// min (atual, bestserver)
	}
}
