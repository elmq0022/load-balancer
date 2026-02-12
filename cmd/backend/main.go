package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var isHealthy bool = true

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from backend"))
	})
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthy := map[string]bool{"healthy": isHealthy}
		data, err := json.Marshal(healthy)
		if err != nil {
			log.Printf("failed to marshall data %v: %v", healthy, err)
			data = []byte("failed to marshall data")
		}

		if isHealthy && err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(data)
	})
	mux.HandleFunc("GET /api/toggle-health", func(w http.ResponseWriter, r *http.Request) {
		isHealthy = !isHealthy
	})

	http.ListenAndServe(":9090", mux)
}
