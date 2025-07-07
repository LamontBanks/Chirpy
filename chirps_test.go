package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestPostChirp(t *testing.T) {
	cfg := initApiConfig()

	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Create new user
	email := "fakeuser@email.com"
	password := "abc123password!"

	_, _, err = createTestUser(cfg, email, password)
	if err != nil {
		t.Error(err)
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

func TestGetChirps(t *testing.T) {
	cfg := initApiConfig()

	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Create 2 users
	email1 := "testuser@gmail.com"
	password1 := "abc123"

	email2 := "testuser2@gmail.com"
	password2 := "xyz890"

	user1, _, err := createTestUser(cfg, email1, password1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user2, _, err := createTestUser(cfg, email2, password2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Post chirps for both users
	messages := []string{"hello, world", "I like turtles", "going on vacation"}

	_, err = postChirp(cfg, user1.Email, password1, messages[0])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, user2.Email, password2, messages[1])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, user2.Email, password2, messages[2])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Get all chirps
	getChirpsReq := httptest.NewRequest("GET", "/api/chirps", nil)
	w := httptest.NewRecorder()
	cfg.getChirps()(w, getChirpsReq)

	chirps := []Chirp{}

	decoder := json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&chirps)
	if err != nil {
		t.Error(err)
	}

	// Assertions
	// All chirps pulled
	if len(chirps) != 3 {
		t.Errorf("failed to GET all chirps: %v", chirps)
	}

	// All messages listed
	for _, message := range messages {
		if !slices.ContainsFunc(chirps, func(c Chirp) bool {
			return c.Body == message
		}) {
			t.Errorf("missing some messages %v, Chirps: %v", messages, chirps)
		}
	}
}

// Logs in the user, posts the given chirp.
// Helper methodfor unit tests
func postChirp(cfg *apiConfig, email, password, chirp string) (Chirp, error) {
	// Log in the user
	loginUserBody := fmt.Sprintf(`{"email": "%v","password": "%v"}`, email, password)
	request := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginUserBody))
	w := httptest.NewRecorder()

	cfg.handlerLogin()(w, request)

	loggedInUser := &LoginResponse{}
	decoder := json.NewDecoder(w.Result().Body)
	err := decoder.Decode(&loggedInUser)
	if err != nil {
		return Chirp{}, err
	}

	// Post chirp
	chirpBody := fmt.Sprintf(`{"body": "%v"}`, chirp)
	chirpRequest := httptest.NewRequest("POST", "/api/chirp", strings.NewReader(chirpBody))
	chirpRequest.Header.Add("Authorization", "Bearer "+loggedInUser.Token)

	// Check response
	w = httptest.NewRecorder()
	cfg.postChirpHandler()(w, chirpRequest)

	if w.Result().StatusCode != http.StatusCreated {
		return Chirp{}, fmt.Errorf("POST /api/chirp/ with auth token failed, response: %v, expected: %v", w.Result().StatusCode, http.StatusCreated)
	}

	// Decode Chirp
	chirpResp := &Chirp{}
	decoder = json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&chirpResp)
	if err != nil {
		return Chirp{}, err
	}

	return *chirpResp, nil
}
