package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) deleteUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CRITICAL: Delete All Users only allowed in dev environment
		if cfg.platform != "dev" {
			sendErrorJSONResponse(w, "cannot DELETE in non-dev environment", http.StatusForbidden, fmt.Errorf("attempted to DELETE in non-dev region"))
			return
		}

		// Count users before delete
		numUsers, err := cfg.db.CountUsers(r.Context())
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Delete users, log number deleted
		err = cfg.db.DeleteUsers(r.Context())
		if err != nil {
			sendErrorJSONResponse(w, "Something went wrong", http.StatusInternalServerError, err)
			return
		}

		// Response
		resp := struct {
			Result string `json:"result"`
		}{
			Result: fmt.Sprintf("All %v users deleted", numUsers),
		}

		SendJSONResponse(w, http.StatusOK, resp)
	}
}
