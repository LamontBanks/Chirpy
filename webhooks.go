package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"` // "user.upgraded"
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerUserUpgraded() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request to validate body parameters
		req := PolkaWebhookRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			sendResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Validate the request fields
		if req.Event != "user.upgraded" {
			sendResponse(w, http.StatusNoContent, fmt.Sprintf("Chirpy Red upgrade handler ignoring event: %v", req.Event))
			return
		}

		if req.Data.UserID == uuid.Nil {
			sendResponse(w, http.StatusNotFound, fmt.Sprintf("Chirpy Red upgrade request missing data.user_id: %v", req))
			return
		}

		// Upgrade the user to "Chirpy Red"
		user, err := cfg.db.UpgradeUserToChirpyRed(r.Context(), req.Data.UserID)
		if err != nil {
			sendResponse(w, http.StatusNotFound, fmt.Sprintf("Chirpy Red upgrade failed for user %v: %v", req.Data.UserID, err))
			return
		}

		// Response
		sendResponse(w, http.StatusNoContent, fmt.Sprintf("Chirpy Red status for user %v changed to %v", user.ID, user.IsChirpyRed))
	}
}
