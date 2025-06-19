package main

import "net/http"

func main() {
	// Matches incoming URL requests to registered patterns and calls the attached handlers
	mux := http.NewServeMux()

	// Handlers
	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
