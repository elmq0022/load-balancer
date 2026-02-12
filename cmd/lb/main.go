package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var targetUrls []*url.URL = []*url.URL{
	&url.URL{
		Scheme: "http",
		Host:   "backend:9090",
	},
}

func main() {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetUrls[0])
		},
	}

	err := http.ListenAndServe(":8080", proxy)
	if err != nil {
		log.Fatal("failed to start server")
	}
}
