package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

var isHealthy bool = true

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request received", "method", r.Method, "path", r.URL.Path)
		w.Write([]byte("hello from backend"))
	})
	mux.HandleFunc("GET /api/healthz/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthy := map[string]bool{"healthy": isHealthy}
		data, err := json.Marshal(healthy)
		if err != nil {
			slog.Error("failed to marshal data", "data", healthy, "error", err)
			data = []byte("failed to marshall data")
		}

		if isHealthy && err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		slog.Info("health check", "healthy", isHealthy, "status", w)
		w.Write(data)
	})
	mux.HandleFunc("GET /api/toggle-health/", func(w http.ResponseWriter, r *http.Request) {
		isHealthy = !isHealthy
		slog.Info("health toggled", "healthy", isHealthy)
	})

	slog.Info("server starting", "port", 9090)
	log.Fatal(http.ListenAndServe(":9090", mux))
}
