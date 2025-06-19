package main

import (
	"net/http"
)

func main() {
	// Matches incoming URL requests to registered patterns and calls the attached handlers
	mux := http.NewServeMux()

	// Handlers
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthHandler)

	// Set handlers, port
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
