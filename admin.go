package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) deleteUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CRITICAL: Delete Users only allowed in dev region
		if cfg.platform != "dev" {
			sendErrorResponse(w, "Cannot DELETE in non-dev environment", http.StatusForbidden, fmt.Errorf("Attempted to DELETE in non-dev region"))
			return
		}

		// Count users before delete
		numUsers, err := cfg.db.CountUsers(r.Context())
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Delete users, log number deleted
		err = cfg.db.DeleteUsers(r.Context())
		if err != nil {
			sendErrorResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}
		sendSuccessResponse(w, "All users deleted", http.StatusOK, fmt.Sprintf("DELETED %v users", int(numUsers)))
	}
}
