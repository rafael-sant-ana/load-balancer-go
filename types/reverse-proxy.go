package types

import (
	"net/http"
)

type ReverseProxy struct {
	Server *ServerInfo
}

func (r *ReverseProxy) SendRequest(req *http.Request) *http.Response {
	var resp *http.Response
	var err error

	fullpath := r.Server.Info.Url + req.URL.Path

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
		panic(err)
	}

	proxyReq.Host = r.Server.Info.Url

	client := &http.Client{}
	resp, err = client.Do(proxyReq)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	return resp
}
