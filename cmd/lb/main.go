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
		Host:   "backend6:6060",
	},
	{
		Scheme: "http",
		Host:   "backend7:7070",
	},
	{
		Scheme: "http",
		Host:   "backend8:8080",
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

	err := http.ListenAndServe(":9090", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: maybe use atomic instead?
		mu.Lock()
		if len(targetUrls) == 0 {
			mu.Unlock()
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		mu.Unlock()
		proxy.ServeHTTP(w, r)
	}))
	if err != nil {
		log.Fatal("failed to start server")
	}
}
