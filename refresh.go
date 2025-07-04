package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/LamontBanks/Chirpy/internal/database"
)

// Returns a new refresh token, if the given refresh token has not expired or been revoked
func (cfg *apiConfig) handlerRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		// Get token from header
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorResponse(w, "Invalid Bearer Token", http.StatusUnauthorized, err)
			return
		}

		// Check if it exists, is not revoked, or is not expired
		tokenInfo, err := cfg.db.GetRefreshTokenInfo(r.Context(), refreshToken)

		// 1. Token doesn't exist
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Invalid bearer token", http.StatusUnauthorized, fmt.Errorf("refresh token %v doesn't exist: %v", refreshToken, err))
			return
		}
		if err != nil {
			sendErrorResponse(w, "", http.StatusUnauthorized, err)
			return
		}

		// 2. Token is not revoked
		// See sql.NullTime for checking NULL SQL values
		// https://pkg.go.dev/database/sql#NullTime
		// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.24.4:src/database/sql/sql.go;l=394
		if t, _ := tokenInfo.RevokedAt.Value(); t != nil {
			sendErrorResponse(w, "Revoked bearer token", http.StatusUnauthorized, err)
			return
		}

		// 3. Token is not expired
		if tokenInfo.ExpiresAt.Before(time.Now()) {
			sendErrorResponse(w, "Expired bearer token", http.StatusUnauthorized, err)
			return
		}

		// Create 60 day refresh token, save to database
		// Create, save, and return a new refresh token for the user
		newToken, err := auth.MakeRefreshToken()
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		refreshTokenDuration, err := time.ParseDuration(auth.REFRESH_TOKEN_DURATION)
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}
		err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     newToken,
			UserID:    tokenInfo.UserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(refreshTokenDuration),
		})
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, 200, response{
			Token: newToken,
		})
	}
}
