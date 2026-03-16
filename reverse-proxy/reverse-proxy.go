package main

import (
	"fmt"
	"io"
	"net/http"
)

var hostname = "pudim.com.br"

func proxy(w http.ResponseWriter, req *http.Request) {
	var resp *http.Response
	var err error

	fullpath := "https://" + hostname + req.URL.Path

	if req.URL.RawQuery != "" {
		fullpath += "?" + req.URL.RawQuery
	}

	proxyReq, err := http.NewRequest(req.Method, fullpath, req.Body)
	for name, headers := range req.Header {
		for _, h := range headers {
			proxyReq.Header.Set(name, h)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	proxyReq.Host = hostname

	client := &http.Client{}
	resp, err = client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}

	defer resp.Body.Close()
	for name, headers := range resp.Header {
		for _, h := range headers {
			w.Header().Set(name, h)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	handler := http.HandlerFunc(proxy)

	fmt.Printf("Server is running\n")
	http.ListenAndServe(":8090", handler)
}
