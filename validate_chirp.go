package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Request, response
	req := struct {
		Body string `json:"body"`
	}{}

	resp := struct {
		Body string `json:"body"`
	}{}

	// Decode request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)

	if err != nil {
		sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
		return
	}

	if len(req.Body) == 0 {
		sendErrorJSONResponse(w, "Something went wrong", http.StatusBadRequest, nil)
		return
	}

	// "Chirps" must be 140 characters or fewer
	if len(req.Body) <= 140 {
		resp.Body = censoredBannedWords(req.Body)
		SendJSONResponse(w, http.StatusOK, resp)
		return
	} else {
		sendErrorJSONResponse(w, "Chirp is too long", http.StatusBadRequest, nil)
		return
	}
}

// Masks banned words from the text, case-insensitive
// Banned words with punctuation, other letters, etc. are skipped
// Example: "Sharbert!" is *not* masked
func censoredBannedWords(originalText string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	censoredText := originalText
	censorMasking := "****"

	for _, bannedWord := range bannedWords {
		splitText := strings.Split(censoredText, " ")

		for i := range splitText {
			if strings.ToLower(splitText[i]) == bannedWord {
				splitText[i] = censorMasking
			}
		}

		// Recombine the string for the next banned word check
		censoredText = strings.Join(splitText, " ")
	}

	return censoredText
}
