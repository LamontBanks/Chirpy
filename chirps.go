package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) postChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type postChirpRequest struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}

		type postChirpResponse struct {
			ChirpID   uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}

		reqBody := postChirpRequest{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		chirpText := reqBody.Body
		userID := reqBody.UserID

		// TODO: Validate userId

		// Validate chirp
		if len(chirpText) == 0 {
			sendErrorResponse(w, "Chirp cannot be empty", http.StatusBadRequest, nil)
		}

		if len(chirpText) > 140 {
			sendErrorResponse(w, "Chirp is too long", http.StatusBadRequest, nil)
		}

		chirpText = censoredBannedWords(chirpText)

		// TODO: Create chirp in database

		// Response
		SendJSONResponse(w, http.StatusCreated, postChirpResponse{
			ChirpID:   uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
			Body:      chirpText,
		})

	}
}
