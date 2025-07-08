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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// Receives text from the users, saves it, and returns the a chirp
func (cfg *apiConfig) postChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Body string `json:"body"`
		}{}

		// Validate Authorization Token
		// 1. Token exists
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		// 2. Token is valid (not expired, etc.)
		userIDFromToken, err := auth.ValidateToken(token, cfg.jwtSecret)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		// 3. Token is associated with a registered user
		_, err = cfg.db.GetUser(r.Context(), userIDFromToken)
		if err == sql.ErrNoRows {
			sendErrorJSONResponse(w, "Invalid User", http.StatusBadRequest, fmt.Errorf("invalid user %v", userIDFromToken))
			return
		}
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Decode request, validate body
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Validate submitted text
		chirpText := req.Body
		if len(chirpText) == 0 {
			sendErrorJSONResponse(w, "Chirp cannot be empty", http.StatusBadRequest, nil)
			return
		}
		if len(chirpText) > 140 {
			sendErrorJSONResponse(w, "Chirp is too long", http.StatusBadRequest, nil)
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
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, http.StatusCreated, Chirp{
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
		chirps := []database.Chirp{}

		// Optional Query Parameter - get chirps by given author_id
		author_id := r.URL.Query().Get("author_id")

		if author_id != "" { // By author_id
			author_uuid, err := uuid.Parse(author_id)
			if err != nil {
				sendErrorJSONResponse(w, "author not found", http.StatusNotFound, err)
				return
			}

			c, err := cfg.db.GetChirpsByUserID(r.Context(), author_uuid)
			if err != nil {
				sendErrorJSONResponse(w, fmt.Sprintf("Failed to get chirps for author %v", author_uuid), http.StatusInternalServerError, err)
				return
			}

			chirps = c
		} else { // All chirps
			c, err := cfg.db.GetChirps(r.Context())
			if err != nil {
				sendErrorJSONResponse(w, "Failed to get all chirps", http.StatusInternalServerError, err)
				return
			}

			chirps = c
		}

		response := []Chirp{}
		for _, c := range chirps {
			response = append(response, Chirp{
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
			sendErrorJSONResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}

		foundChirp, err := cfg.db.GetChirpByID(r.Context(), id)
		if err != nil {
			sendErrorJSONResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}

		SendJSONResponse(w, http.StatusOK, Chirp{
			ID:        foundChirp.ID,
			CreatedAt: foundChirp.CreatedAt,
			UpdatedAt: foundChirp.UpdatedAt,
			Body:      foundChirp.Body,
			UserID:    foundChirp.UserID,
		})
	}
}

func (cfg *apiConfig) deleteChirpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the chirpID, check if it exists
		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			sendErrorJSONResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}

		chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
		if err == sql.ErrNoRows {
			sendErrorJSONResponse(w, "Chirp not found", http.StatusNotFound, err)
			return
		}
		if err != nil && err != sql.ErrNoRows {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Get userID from auth token
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		userIDFromToken, err := auth.ValidateToken(token, cfg.jwtSecret)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		// Verify the chirp was made by the user
		if userIDFromToken != chirp.UserID {
			sendResponse(w, http.StatusForbidden, fmt.Sprintf("user %v tried deleting unowned chirp %v", userIDFromToken, chirp))
			return
		}

		// Delete Chirp
		deletedChirp, err := cfg.db.DeleteChirpByID(r.Context(), chirpID)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		sendResponse(w, http.StatusNoContent, fmt.Sprintf("user %v deleted chirp %v", userIDFromToken, deletedChirp.ID))
	}
}
