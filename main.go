package main

import (
	"encoding/json"
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Expected params for the request JSON body
	type requestBody struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		Valid bool `json:"valid"`
	}

	// Decode into the params struct
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)

	if err != nil {
		sendErrorResponse(w, "Something went wrong", 500, err)
		return
	}

	if len(request.Body) == 0 {
		sendErrorResponse(w, "Something went wrong", 400, nil)
		return
	}

	// "Chirps" must be 140 characters or fewer
	if len(request.Body) <= 140 {
		success := successResponse{
			Valid: true,
		}

		sendJSONResponse(w, 200, success)
		return
	} else {
		sendErrorResponse(w, "Chirp is too long", 400, nil)
		return
	}
}

// Display metrics (ex: page count)
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	responseHTML := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(responseHTML))
}

// Resets the hit count to 0
func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
