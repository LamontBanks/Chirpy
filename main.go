package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/LamontBanks/Chirpy/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// Store stateful data between API calls
type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	godotenv.Load() // .env at root

	// Initialize database
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic("Error connecting to the database")
	}

	dbQueries := database.New(db)

	// Determine execution platform (dev, qa, prod)
	platform := os.Getenv("PLATFORM")

	// Set values into config
	cfg := &apiConfig{
		db:       dbQueries,
		platform: platform,
	}
	cfg.fileServerHits.Store(0)

	// Set Endpoints
	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	mux.HandleFunc("POST /api/users", cfg.createUserHandler())

	// Start server
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

// Increments the number of hits to the server for the given endpoint handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Display metrics (ex: page count)
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	responseHTML := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseHTML))
}

// Resets the hit count to 0
func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
