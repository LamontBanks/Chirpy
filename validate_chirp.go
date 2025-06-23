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
		Body string `json:"body"`
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
		resp.Body = censoredBannedWords(req.Body)
		sendJSONResponse(w, 200, resp)
		return
	} else {
		sendErrorResponse(w, "Chirp is too long", 400, nil)
		return
	}
}

func censoredBannedWords(originalChirp string) string {
	censoredText := originalChirp
	censorMasking := "****"

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range bannedWords {
		censoredText = strings.ReplaceAll(censoredText, word, censorMasking)
	}

	return censoredText
}
