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
	setup()
	defer tearDown()

	cases := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "New User 1",
			email:    "fakeuser@email.com",
			password: "abc!password123",
		},
		{
			name:     "New User 2",
			email:    "fakeuser2@email.com",
			password: "abc123",
		},
	}

	cfg := initApiConfig()
	for _, c := range cases {
		input := fmt.Sprintf(`{"email": "%v", "password": "%v"}`, c.email, c.password)

		newUser, responseCode, err := createTestUser(cfg, c.email, c.password)
		if err != nil {
			t.Errorf("failed to create user %v: %v", input, err)
			t.FailNow()
		}

		// Validate user fields
		assertEquals(responseCode, http.StatusCreated, newUser, t)
		assertEquals(newUser.Email, c.email, newUser, t)
		assertEquals(newUser.IsChirpyRed, false, newUser, t)
		assertEquals(newUser.CreatedAt.IsZero(), false, newUser, t)
		assertEquals(newUser.UpdatedAt.IsZero(), false, newUser, t)
		assertEquals(uuid.Validate(newUser.ID.String()), nil, newUser, t)
	}
}

func deleteAllUsersAndPosts(cfg *apiConfig) error {
	if cfg.platform != "dev" {
		return fmt.Errorf("cannot call /api/reset in non-dev environment")
	}

	resetRequest := httptest.NewRequest("POST", "/api/reset", nil)
	w := httptest.NewRecorder()
	cfg.deleteUsersHandler()(w, resetRequest)

	if w.Result().StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete all users, posts")
	}

	return nil
}

// Attempts to create a new user using the JSON string
// Returns the status code and any error
func createTestUser(cfg *apiConfig, email, password string) (*User, int, error) {
	request := httptest.NewRequest("POST", "/api/users", strings.NewReader(fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password)))
	w := httptest.NewRecorder()
	cfg.createUserHandler()(w, request)

	response := w.Result()

	if response.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return &User{}, response.StatusCode, fmt.Errorf("error creating new user: %v", err)
		}
		return &User{}, response.StatusCode, fmt.Errorf("error creating new user: %v", string(body))
	}

	user := User{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&user)

	return &user, w.Result().StatusCode, err
}

// Create user, also returns passwords - indices correspond to each user
// Unit test helper function
func createTestUsers(cfg *apiConfig, numUsers int) ([]User, []string, error) {
	newUsers := []User{}
	passwords := []string{}

	for i := range numUsers {
		email := fmt.Sprintf("testuser_%v@gmail.com", i)
		pw := fmt.Sprintf("abc00%v", i)

		user, _, err := createTestUser(cfg, email, pw)
		if err != nil {
			return []User{}, []string{}, err
		}
		passwords = append(passwords, pw)

		newUsers = append(newUsers, *user)
	}

	return newUsers, passwords, nil
}
