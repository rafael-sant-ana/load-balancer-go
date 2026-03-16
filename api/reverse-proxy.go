package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var src_page = "https://pokeapi.co/api/v2"

func proxy(w http.ResponseWriter, req *http.Request) {
	route := req.PathValue("path")
	for name, headers := range req.Header {
		for _, h := range headers {
			req.Header.Set(name, h)
		}
	}

	var resp *http.Response
	var err error

	fullpath := strings.Join([]string{src_page, route}, "/")

	fmt.Println(fullpath)
	switch req.Method { // depois tem que deixar dinamico pra custom methods
	case "GET":
		resp, err = http.Get(fullpath)

		if err != nil {
			fmt.Println(err)
			// retornar na response
		}

	case "POST":
		resp, err = http.Post(fullpath, req.Header.Get("Content-Type"), req.Body)

		if err != nil {
			fmt.Println(err)
		}

	case "DELETE":

	case "PUT":

	case "PATCH":

	}
	defer resp.Body.Close()
	for name, headers := range resp.Header {
		for _, h := range headers {
			w.Header().Set(name, h)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// msg
	}

	w.Write(body)
}

func main() {
	http.HandleFunc("/{path...}", proxy)

	fmt.Printf("Server is running\n")
	http.ListenAndServe(":8090", nil)
}
