package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUserCreation(t *testing.T) {
	cfg := initApiConfig()

	// Delete users from database
	deleteAllUsersAndPosts(cfg, t)

	// Create new user
	input := `{"email": "fakeuser@email.com", "password": "abc123password!"}`
	newUser, responseCode, err := createTestUser(cfg, input)
	if err != nil {
		t.Errorf("failed to create user %v: %v", input, err)
		t.FailNow()
	}

	// Validate user fields
	assertEqual(responseCode, http.StatusCreated, input, t)
	assertEqual(newUser.Email, "fakeuser@email.com", input, t)

	if err := uuid.Validate(newUser.ID.String()); err != nil {
		t.Errorf("id is not a valid UUID: %v", newUser)
		t.FailNow()
	}

	// TODO: Verify timestamps
}

func deleteAllUsersAndPosts(cfg *apiConfig, t *testing.T) {
	if cfg.platform != "dev" {
		t.Errorf("cannot call /api/reset in non-dev environment")
		t.FailNow()
	}

	resetRequest := httptest.NewRequest("POST", "/api/reset", nil)
	w := httptest.NewRecorder()
	cfg.deleteUsersHandler()(w, resetRequest)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("failed to delete all users: %v", w.Result())
		t.FailNow()
	}
}

// Attempts to create a new user using the JSON string
// Returns the status code and any error
func createTestUser(cfg *apiConfig, inputBody string) (*User, int, error) {
	// New user request
	request := httptest.NewRequest("POST", "/api/users", strings.NewReader(inputBody))
	w := httptest.NewRecorder()

	// Make request
	cfg.createUserHandler()(w, request)

	// Get response
	response := w.Result()

	// Return error if user creation fails
	if response.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return &User{}, response.StatusCode, fmt.Errorf("error creating new user: %v", err)
		}
		return &User{}, response.StatusCode, fmt.Errorf("error creating new user: %v", string(body))
	}

	// Otheriwse, decode response into User struct, return
	user := User{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&user)

	return &user, w.Result().StatusCode, err
}
