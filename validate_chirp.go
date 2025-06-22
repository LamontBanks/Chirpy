package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Request, response
	type requestBody struct {
		Body string `json:"body"`
	}
	req := requestBody{}

	type responseBody struct {
		Valid bool `json:"valid"`
	}
	resp := responseBody{}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	// Validate chirp
	if err != nil {
		sendErrorResponse(w, "Something went wrong", 500, err)
		return
	}

	if len(req.Body) == 0 {
		sendErrorResponse(w, "Something went wrong", 400, nil)
		return
	}

	// "Chirps" must be 140 characters or fewer
	if len(req.Body) <= 140 {
		resp.Valid = true
		sendJSONResponse(w, 200, resp)
		return
	} else {
		sendErrorResponse(w, "Chirp is too long", 400, nil)
		return
	}
}

func censoredBannedWords(w http.ResponseWriter, r *http.Request) {
	// Request, response
	type request struct {
		Body string
	}
	req := request{}

	type response struct {
		Body string
	}
	resp := response{}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Something went wrong", 500, err)
	}

	// Replace banned words
	censored_chirp := req.Body
	censoredText := "****"

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range bannedWords {
		censored_chirp = strings.ReplaceAll(censored_chirp, word, censoredText)
	}

	// Response
	resp.Body = censored_chirp
	sendJSONResponse(w, 200, resp)
}
