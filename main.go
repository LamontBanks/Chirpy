package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Store stateful data between API calls
type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	cfg := &apiConfig{}
	cfg.fileServerHits.Store(0)

	// Matches incoming URL requests to registered patterns and calls the attached handlers
	mux := http.NewServeMux()

	// Handlers
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /reset", cfg.resetMetricsHandler)

	// Set handlers, port
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

// Handlers

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// Increments the number of hits to the server for the given endpoint handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	fmt.Fprintf(w, "Hits: %v", cfg.fileServerHits.Load())
}

// Resets the hit count to 0
func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
