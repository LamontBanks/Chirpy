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
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Check for required elements
		if reqBody.Email == "" {
			sendErrorJSONResponse(w, "Email required", http.StatusBadRequest, nil)
			return
		}

		if reqBody.Password == "" {
			sendErrorJSONResponse(w, "Password required", http.StatusBadRequest, nil)
			return
		}

		// Save password
		hashedPassword, err := auth.HashPassword(reqBody.Password)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid Password", http.StatusBadRequest, err)
			return
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
			sendErrorJSONResponse(w, msg, 500, err)
			return
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

// Updates the user's email and password, based on the provided authentication token
func (cfg *apiConfig) updateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Read userID from the auth token
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid bearer token", http.StatusUnauthorized, err)
			return
		}

		userID, err := auth.ValidateToken(token, cfg.jwtSecret)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid bearer token", http.StatusUnauthorized, err)
			return
		}

		// Decode request, check new email and password values
		req := request{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&req)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		if req.Email == "" {
			sendErrorJSONResponse(w, "Email required", http.StatusBadRequest, fmt.Errorf("empty email: %v", req.Email))
			return
		}
		if req.Password == "" {
			sendErrorJSONResponse(w, "Password required", http.StatusBadRequest, fmt.Errorf("empty password: %v", req.Password))
			return
		}

		new_hashed_password, err := auth.HashPassword(req.Password)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Update user info
		updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
			ID:             userID,
			Email:          req.Email,
			HashedPassword: new_hashed_password,
			UpdatedAt:      time.Now(),
		})
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, http.StatusOK, User{
			ID:        updatedUser.ID,
			Email:     updatedUser.Email,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
		})

	}
}
