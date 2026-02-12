package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var isHealthy bool = true
var port string = os.Getenv("PORT")

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request received", "method", r.Method, "path", r.URL.Path)
		w.Write([]byte("hello from backend on port: " + port))
	})
	mux.HandleFunc("GET /api/healthz/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthy := map[string]bool{"healthy": isHealthy}
		data, err := json.Marshal(healthy)
		if err != nil {
			slog.Error("failed to marshal data", "data", healthy, "error", err)
			data = []byte("failed to marshall data")
		}

		statusCode := http.StatusInternalServerError
		if isHealthy && err == nil {
			statusCode = http.StatusOK
		}
		w.WriteHeader(statusCode)

		slog.Info("health check", "healthy", isHealthy, "status", statusCode)
		w.Write(data)
	})
	mux.HandleFunc("GET /api/toggle-health/", func(w http.ResponseWriter, r *http.Request) {
		isHealthy = !isHealthy
		slog.Info("health toggled for backend port: "+port, "healthy", isHealthy)
	})

	slog.Info("server starting", "port", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
