package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var targetUrl *url.URL = &url.URL{
	Scheme: "http",
	Host:   "backend:9090",
}

func main() {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetUrl)
		},
	}

	err := http.ListenAndServe(":8080", proxy)
	if err != nil {
		log.Fatal("failed to start server")
	}
}
