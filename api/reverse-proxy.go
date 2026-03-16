package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			req.Header.Set(name, h)
		}
	}

	var resp *http.Response
	var err error

	switch req.Method { // depois tem que deixar dinamico pra custom methods
	case "GET":
		resp, err = http.Get("https://google.com.br")

		if err != nil {
			fmt.Println(err)
		}

	case "POST":

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

	response_text := ""
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		response_text += scanner.Text()
	}

	w.Write([]byte(response_text))
}

func main() {
	http.HandleFunc("/hello", hello)
	// http.HandleFunc("")

	fmt.Printf("Server is running\n")
	http.ListenAndServe(":8090", nil)
}
