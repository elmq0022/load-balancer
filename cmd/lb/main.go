package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var targetUrls []*url.URL = []*url.URL{
	{
		Scheme: "http",
		Host:   "backend:9090",
	},
	{
		Scheme: "http",
		Host:   "backend:9090",
	},
}

var server int = 0
var mu sync.Mutex

func main() {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			mu.Lock()
			r.SetURL(targetUrls[server])
			server = (server + 1) % len(targetUrls)
			mu.Unlock()
		},
	}

	err := http.ListenAndServe(":8080", proxy)
	if err != nil {
		log.Fatal("failed to start server")
	}
}
