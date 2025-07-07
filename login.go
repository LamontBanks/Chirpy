package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LamontBanks/Chirpy/internal/auth"
	"github.com/LamontBanks/Chirpy/internal/database"
	"github.com/google/uuid"
)

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Password string `json:"password"`
			Email    string `json:"email"`
		}{}

		// Decode request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Validate required fields
		if req.Email == "" {
			sendErrorJSONResponse(w, "Email must not be blank", http.StatusBadRequest, nil)
			return
		}

		if req.Password == "" {
			sendErrorJSONResponse(w, "Password required", http.StatusBadRequest, nil)
			return
		}

		// Check user password
		user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			sendErrorJSONResponse(w, "incorrect email or password", http.StatusUnauthorized, err)
			return
		}

		err = auth.CheckPasswordHash(req.Password, user.HashedPassword)
		if err != nil {
			sendErrorJSONResponse(w, "incorrect email or password", http.StatusUnauthorized, err)
			return
		}

		// Check JWT hasn't expired
		tokenDuration, err := time.ParseDuration(auth.JWT_TOKEN_DURATION)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, tokenDuration)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Create 60 day refresh token, save to database
		refreshTokenDuration, err := time.ParseDuration(auth.REFRESH_TOKEN_DURATION)
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(refreshTokenDuration),
		})
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		SendJSONResponse(w, 200, LoginResponse{
			ID:           user.ID,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			IsChirpyRed:  user.IsChirpyRed,
			Token:        token,
			RefreshToken: refreshToken,
		})
	}
}
