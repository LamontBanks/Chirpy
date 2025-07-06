package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/LamontBanks/Chirpy/internal/database"
)

// Returns a new JWT token, if the given refresh token is still valid
func (cfg *apiConfig) handlerRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		// Get token from header
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		// Check if it exists, is not revoked, or is not expired
		refreshTokenInfo, err := cfg.db.GetRefreshTokenInfo(r.Context(), refreshToken)

		// 1. Refresh Token doesn't exist
		if err == sql.ErrNoRows {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, fmt.Errorf("refresh token %v doesn't exist: %v", refreshToken, err))
			return
		}
		if err != nil {
			sendErrorJSONResponse(w, "", http.StatusUnauthorized, err)
			return
		}

		// 2. Refresh Token is not revoked
		// See sql.NullTime for checking NULL SQL values
		// https://pkg.go.dev/database/sql#NullTime
		// Source: https://cs.opensource.google/go/go/+/refs/tags/go1.24.4:src/database/sql/sql.go;l=394
		if t, _ := refreshTokenInfo.RevokedAt.Value(); t != nil {
			sendErrorJSONResponse(w, "Revoked bearer token", http.StatusUnauthorized, err)
			return
		}

		// 3. Refresh Token is not expired
		if refreshTokenInfo.ExpiresAt.Before(time.Now()) {
			sendErrorJSONResponse(w, "Expired bearer token", http.StatusUnauthorized, err)
			return
		}

		// Create, send new JWT token
		jwtTokenDuration, _ := time.ParseDuration(auth.JWT_TOKEN_DURATION)
		newJWTToken, err := auth.MakeJWT(refreshTokenInfo.UserID, cfg.jwtSecret, jwtTokenDuration)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, 200, response{
			Token: newJWTToken,
		})
	}
}

// Revokes the refresh token in the header
func (cfg *apiConfig) handlerRevoke() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, err)
			return
		}

		// Check if refresh token exists
		existingRefreshTokenInfo, err := cfg.db.GetRefreshTokenInfo(r.Context(), refreshToken)
		if err == sql.ErrNoRows {
			sendErrorJSONResponse(w, "Invalid User", http.StatusUnauthorized, fmt.Errorf("refresh token %v doesn't exist: %v", refreshToken, err))
			return
		}
		if err != nil {
			sendResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Token exists, set RevokeAt time
		// RevokedAt can be NULL, so need to wrap the new Time in sql.NullTime type
		revokeTime := &sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		err = cfg.db.SetRefreshTokenRevokeAtTime(r.Context(), database.SetRefreshTokenRevokeAtTimeParams{
			Token:     existingRefreshTokenInfo.Token,
			RevokedAt: *revokeTime,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			sendResponse(w, http.StatusInternalServerError, err.Error())
		}

		// Response
		sendResponse(w, http.StatusNoContent, "")
	}
}
