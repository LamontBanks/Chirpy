package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/LamontBanks/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Token     string    `json:"token"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// Receives text from the users, saves it, and returns the a chirp
func (cfg *apiConfig) postChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Body string `json:"body"`
		}

		// Validate Authorization Token
		// 1. Token exists
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorResponse(w, "Invalid JWT", http.StatusUnauthorized, err)
			return
		}

		// 2. Token is valid (not expired, etc.)
		userIDFromToken, err := auth.ValidateToken(token, cfg.jwtSecret)
		if err != nil {
			sendErrorResponse(w, "Invalid JWT", http.StatusUnauthorized, err)
			return
		}

		// 3. Token is associated with a registered user
		_, err = cfg.db.GetUser(r.Context(), userIDFromToken)
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Invalid User", http.StatusBadRequest, fmt.Errorf("invalid user %v", userIDFromToken))
			return
		}
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Decode request to validate body parameters
		reqBody := request{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&reqBody)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Validate submitted text
		chirpText := reqBody.Body
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
			UserID:    userIDFromToken,
			Body:      chirpText,
		})
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, http.StatusCreated, chirp{
			ID:        savedChirp.ID,
			CreatedAt: savedChirp.CreatedAt,
			UpdatedAt: savedChirp.UpdatedAt,
			UserID:    savedChirp.UserID,
			Body:      savedChirp.Body,
		})
	}
}

func (cfg *apiConfig) getChirps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := []chirp{}

		// Gets in ascendaing order by created_at
		allChirps, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			sendErrorResponse(w, "Failed to gets all chirps", http.StatusInternalServerError, err)
			return
		}

		for _, c := range allChirps {
			response = append(response, chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}

		SendJSONResponse(w, http.StatusOK, response)
	}
}

func (cfg *apiConfig) getChirpByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			sendErrorResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}

		foundChirp, err := cfg.db.GetChirpByID(r.Context(), id)
		if err != nil {
			sendErrorResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}

		SendJSONResponse(w, http.StatusOK, chirp{
			ID:        foundChirp.ID,
			CreatedAt: foundChirp.CreatedAt,
			UpdatedAt: foundChirp.UpdatedAt,
			Body:      foundChirp.Body,
			UserID:    foundChirp.UserID,
		})
	}
}
