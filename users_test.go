package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUserCreation(t *testing.T) {
	cfg := initApiConfig()

	// Delete users from database
	deleteUsersResponseCode, err := deleteAllUsers(cfg)
	if deleteUsersResponseCode != http.StatusOK {
		t.Errorf("failed to delete all users %v", err)
	}

	// Create new user
	input := `{"email": "fakeuser@email.com","password": "abc123password!"}`
	newUser, responseCode, err := createTestUser(cfg, input)

	if err != nil {
		t.Errorf("failed to create user %v: %v", input, err)
	}

	// Validate user fields
	assertEqual(responseCode, http.StatusCreated, input, t)
	assertEqual(newUser.Email, "fakeuser@email.com", input, t)

	if err := uuid.Validate(newUser.ID.String()); err != nil {
		t.Errorf("id is not a valid UUID: %v", newUser)
	}

	// TODO: Verify timestamps
}

func deleteAllUsers(cfg *apiConfig) (int, error) {
	if cfg.platform != "dev" {
		return http.StatusInternalServerError, fmt.Errorf("cannot call /api/reset - not in dev environment")
	}

	resetRequest := httptest.NewRequest("POST", "/api/reset", nil)
	w := httptest.NewRecorder()

	cfg.deleteUsersHandler()(w, resetRequest)
	if w.Result().StatusCode != http.StatusOK {
		return w.Result().StatusCode, fmt.Errorf("failed to delete users: %v", w.Result())
	}

	return w.Result().StatusCode, nil
}

// Attempts to create a new user using the JSON string
// Returns the status code and any error
func createTestUser(cfg *apiConfig, inputBody string) (User, int, error) {
	request := httptest.NewRequest("POST", "/api/users", strings.NewReader(inputBody))
	w := httptest.NewRecorder()

	cfg.createUserHandler()(w, request)

	response := w.Result()
	decoder := json.NewDecoder(response.Body)

	user := User{}

	err := decoder.Decode(&user)

	return user, w.Result().StatusCode, err
}
