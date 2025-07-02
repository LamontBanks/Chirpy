package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/LamontBanks/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Wrap functions in a closure to get access to the database
func (cfg *apiConfig) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Password string `json:"password"`
			Email    string `json:"email"`
		}

		// Decode request
		reqBody := requestBody{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Check for required elements
		if reqBody.Email == "" {
			sendErrorResponse(w, "Email must not be blank", http.StatusBadRequest, nil)
			return
		}

		if reqBody.Password == "" {
			sendErrorResponse(w, "Password required", http.StatusBadRequest, nil)
			return
		}

		// Save password
		hashedPassword, err := auth.HashPassword(reqBody.Password)
		if err != nil {
			sendErrorResponse(w, "Invalid Password", http.StatusBadRequest, err)
		}

		// Create user
		dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			ID:             uuid.New(),
			Email:          reqBody.Email,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			HashedPassword: hashedPassword,
		})

		if err != nil {
			msg := fmt.Sprintf("Unable to create user with email %v", reqBody.Email)
			sendErrorResponse(w, msg, 500, err)
		}

		// Map from database.User to custom User type
		user := User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
		}

		// Success Response
		SendJSONResponse(w, http.StatusCreated, user)
	}
}
