package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/database"
	"github.com/google/uuid"
)

// Receives text from the users, saves it, and returns the a chirp
func (cfg *apiConfig) postChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request, response format
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

		// Decode request
		reqBody := postChirpRequest{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		chirpText := reqBody.Body
		userID := reqBody.UserID

		// Check the user exists
		_, err = cfg.db.GetUser(r.Context(), userID)
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Invalid User", http.StatusBadRequest, fmt.Errorf("invalid user %v attempted to post", userID))
			return
		}
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Validate submitted text
		if len(chirpText) == 0 {
			sendErrorResponse(w, "Chirp cannot be empty", http.StatusBadRequest, nil)
			return
		}
		if len(chirpText) > 140 {
			sendErrorResponse(w, "Chirp is too long", http.StatusBadRequest, nil)
			return
		}
		chirpText = censoredBannedWords(chirpText)

		// Create chirp in database
		savedChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
			Body:      chirpText,
		})
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, http.StatusCreated, postChirpResponse{
			ChirpID:   savedChirp.ID,
			CreatedAt: savedChirp.CreatedAt,
			UpdatedAt: savedChirp.UpdatedAt,
			UserID:    savedChirp.UserID,
			Body:      savedChirp.Body,
		})
	}
}
