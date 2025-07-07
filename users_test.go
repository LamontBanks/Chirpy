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

	deleteAllUsersAndPosts(cfg, t)

	// Create new user
	email := "fakeuser@email.com"
	password := "abc!password123"
	input := fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password)

	newUser, responseCode, err := createTestUser(cfg, input)
	if err != nil {
		t.Errorf("failed to create user %v: %v", input, err)
	}

	// Validate user fields
	assertEquals(responseCode, http.StatusCreated, newUser, t)
	assertEquals(newUser.Email, email, newUser, t)
	assertEquals(newUser.IsChirpyRed, false, newUser, t)
	assertEquals(newUser.CreatedAt.IsZero(), false, newUser, t)
	assertEquals(newUser.UpdatedAt.IsZero(), false, newUser, t)
	assertEquals(uuid.Validate(newUser.ID.String()), nil, newUser, t)
}

func deleteAllUsersAndPosts(cfg *apiConfig, t *testing.T) {
	if cfg.platform != "dev" {
		t.Errorf("cannot call /api/reset in non-dev environment")
		t.FailNow()
	}

	resetRequest := httptest.NewRequest("POST", "/api/reset", nil)
	w := httptest.NewRecorder()
	cfg.deleteUsersHandler()(w, resetRequest)

	assertEquals(w.Result().StatusCode, http.StatusOK, resetRequest, nil)
}

// Attempts to create a new user using the JSON string
// Returns the status code and any error
func createTestUser(cfg *apiConfig, inputBody string) (*User, int, error) {
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
