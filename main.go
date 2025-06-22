package main

import (
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

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
