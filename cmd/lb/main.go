package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

var configUrls []*url.URL = []*url.URL{
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

var servers atomic.Pointer[[]*url.URL]

var server int = 0
var mu sync.Mutex

func main() {
	servers.Store(&configUrls)
	go func() {
		for range time.Tick(10 * time.Second) {
			var healthy []*url.URL
			for _, u := range configUrls {
				if isHealthy(u.Host) {
					healthy = append(healthy, u)
				}
			}
			servers.Store(&healthy)
		}
	}()

	err := http.ListenAndServe(":9090", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := *servers.Load()
		if len(s) == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		proxy := &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				mu.Lock()
				server = server % len(s)
				r.SetURL(s[server])
				server = (server + 1) % len(s)
				mu.Unlock()
			},
		}
		proxy.ServeHTTP(w, r)
	}))
	if err != nil {
		log.Fatal("failed to start server")
	}
}

func isHealthy(host string) bool {
	req, err := http.NewRequest("GET", "http://"+host+"/api/healthz/", nil)
	if err != nil {
		return false
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
