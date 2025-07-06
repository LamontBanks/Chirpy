package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostChirp(t *testing.T) {
	cfg := initApiConfig()

	deleteAllUsersAndPosts(cfg, t)

	// Create new user
	email := "fakeuser@email.com"
	password := "abc123password!"
	createUserBody := fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password)

	_, _, err := createTestUser(cfg, createUserBody)
	if err != nil {
		t.Errorf("failed to create user %v: %v", createUserBody, err)
		t.FailNow()
	}

	// Log in, get the JWT token required for posting
	loginUserBody := fmt.Sprintf(`{"email": "%v","password": "%v"}`, email, password)
	request := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginUserBody))

	w := httptest.NewRecorder()
	cfg.handlerLogin()(w, request)

	// Read response
	loggedInUser := &LoginResponse{}
	decoder := json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&loggedInUser)
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	// Create chirp, pass the auth token
	chirp := `{"body": "Hello world"}`
	chirpRequest := httptest.NewRequest("POST", "/api/chirp", strings.NewReader(chirp))
	chirpRequest.Header.Add("Authorization", "Bearer "+loggedInUser.Token)

	// Post chirp, check response
	w = httptest.NewRecorder()
	cfg.postChirpHandler()(w, chirpRequest)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("POST /api/chirp/ with auth token failed, response: %v, expected: %v", w.Result().StatusCode, http.StatusCreated)
	}
}
