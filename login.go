package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Password string `json:"password"`
			Email    string `json:"email"`

			// Optional
			ExpiresInSeconds int `json:"expires_in_seconds"`
		}

		// Decode request
		reqBody := requestBody{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Validate required fields
		if reqBody.Email == "" {
			sendErrorResponse(w, "Email must not be blank", http.StatusBadRequest, nil)
			return
		}

		if reqBody.Password == "" {
			sendErrorResponse(w, "Password required", http.StatusBadRequest, nil)
			return
		}

		// Default to 1 hour token expiration; use client-provided expiration if it falls within bounds
		tokenExpirationDuration, err := time.ParseDuration("1h")
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		if reqBody.ExpiresInSeconds > 0 && reqBody.ExpiresInSeconds <= int(tokenExpirationDuration.Seconds()) {
			d := strconv.Itoa(reqBody.ExpiresInSeconds)
			tokenExpirationDuration, err = time.ParseDuration(d + "s")

			if err != nil {
				sendErrorResponse(w, fmt.Sprintf("Invalid expires_in_seconds value: %v", reqBody.ExpiresInSeconds), http.StatusBadRequest, nil)
				return
			}
		}

		// Check user password
		user, err := cfg.db.GetUserByEmail(r.Context(), reqBody.Email)
		if err != nil {
			sendErrorResponse(w, "incorrect email or password", http.StatusUnauthorized, err)
			return
		}

		err = auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
		if err != nil {
			sendErrorResponse(w, "incorrect email or password", http.StatusUnauthorized, err)
			return
		}

		// Create a new JWT token
		token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, tokenExpirationDuration)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, 200, LoginResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		})
	}
}
