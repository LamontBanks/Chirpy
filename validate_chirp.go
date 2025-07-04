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
		sendErrorJSONResponse(w, "Something went wrong", 500, err)
		return
	}

	if len(req.Body) == 0 {
		sendErrorJSONResponse(w, "Something went wrong", 400, nil)
		return
	}

	// "Chirps" must be 140 characters or fewer
	if len(req.Body) <= 140 {
		resp.Body = censoredBannedWords(req.Body)
		SendJSONResponse(w, http.StatusOK, resp)
		return
	} else {
		sendErrorJSONResponse(w, "Chirp is too long", 400, nil)
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
		// Split text by spaces to isolate words
		splitText := strings.Split(censoredText, " ")

		for i := range splitText {
			// Convert to lowercase for case-insensitve comparison
			if strings.ToLower(splitText[i]) == bannedWord {
				splitText[i] = censorMasking
			}
		}

		// Recombine the string for the next banned word check
		censoredText = strings.Join(splitText, " ")
	}

	return censoredText
}
