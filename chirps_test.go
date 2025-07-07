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

	// User 1
	loggedInUser1, err := loginUser(cfg, user1.Email, password1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, loggedInUser1.Token, messages[0])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// User 2
	loggedInUser2, err := loginUser(cfg, user2.Email, password2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, loggedInUser2.Token, messages[1])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, loggedInUser2.Token, messages[2])
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

func TestDeleteChirp(t *testing.T) {
	cfg := initApiConfig()

	err := deleteAllUsersAndPosts(cfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Create user, login
	email := "testuser@gmail.com"
	password := "abc123"

	user, _, err := createTestUser(cfg, email, password)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	loggedInUser, err := loginUser(cfg, user.Email, password)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Post multiple chirps
	chirp1, err := postChirp(cfg, loggedInUser.Token, "hello, world")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, loggedInUser.Token, "I like turtles")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = postChirp(cfg, loggedInUser.Token, "going on vacation")
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

	// Quick check to ensure all the chirps posted
	assertEquals(len(chirps), 3, chirps, t)

	// Delete specific chirps
	deleteChirpRequest := httptest.NewRequest("DELETE", "/api/chirps/", nil)
	deleteChirpRequest.SetPathValue("chirpID", chirp1.ID.String())
	deleteChirpRequest.Header.Add("Authorization", "Bearer "+loggedInUser.Token)
	w = httptest.NewRecorder()

	cfg.deleteChirpHandler()(w, deleteChirpRequest)

	// Success response
	assertEquals(w.Result().StatusCode, http.StatusNoContent, deleteChirpRequest, t)

	// Ensure chirp was removed
	getChirpsReq = httptest.NewRequest("GET", "/api/chirps", nil)
	w = httptest.NewRecorder()
	cfg.getChirps()(w, getChirpsReq)

	chirps = []Chirp{}

	decoder = json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&chirps)
	if err != nil {
		t.Error(err)
	}

	assertEquals(len(chirps), 2, chirps, t)
}

// Helper method to posts the chirp for the user
// User must login to get their auth token
func postChirp(cfg *apiConfig, userAuthToken, chirp string) (Chirp, error) {
	// Post chirp
	chirpBody := fmt.Sprintf(`{"body": "%v"}`, chirp)
	chirpRequest := httptest.NewRequest("POST", "/api/chirp", strings.NewReader(chirpBody))
	chirpRequest.Header.Add("Authorization", "Bearer "+userAuthToken)

	// Check response
	w := httptest.NewRecorder()
	cfg.postChirpHandler()(w, chirpRequest)

	if w.Result().StatusCode != http.StatusCreated {
		return Chirp{}, fmt.Errorf("POST /api/chirp/ with auth token failed, response: %v, expected: %v", w.Result().StatusCode, http.StatusCreated)
	}

	// Decode Chirp
	chirpResp := &Chirp{}
	decoder := json.NewDecoder(w.Result().Body)
	err := decoder.Decode(&chirpResp)
	if err != nil {
		return Chirp{}, err
	}

	return *chirpResp, nil
}
