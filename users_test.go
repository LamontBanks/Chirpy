package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	setup()
	tests := m.Run()
	tearDown()
	os.Exit(tests)
}

func setup() {
	cfg := initApiConfig()
	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func tearDown() {
	cfg := initApiConfig()
	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func TestUserCreation(t *testing.T) {
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

// Create multiple users, unit test helper function
func createMultipleUsers(cfg *apiConfig, numUsers int) ([]User, error) {
	newUsers := []User{}

	for i := range numUsers {
		email := fmt.Sprintf("testuser_%v", i)
		password := fmt.Sprintf("abc00%v", i)

		user, _, err := createTestUser(cfg, email, password)
		if err != nil {
			return []User{}, err
		}

		newUsers = append(newUsers, *user)
	}

	return newUsers, nil
}
