package main

import (
	"encoding/json"
	"net/http"

	"github.com/LamontBanks/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin() http.HandlerFunc {
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

		// Validate request
		if reqBody.Email == "" {
			sendErrorResponse(w, "Email must not be blank", http.StatusBadRequest, nil)
			return
		}

		if reqBody.Password == "" {
			sendErrorResponse(w, "Password required", http.StatusBadRequest, nil)
			return
		}

		// Get hash from db, compare to user-provided password
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

		SendJSONResponse(w, 200, User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
}
