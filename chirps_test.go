package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestPostChirp(t *testing.T) {
	cases := []struct {
		email         string
		password      string
		chirpReqBody  string
		expectedChirp Chirp
	}{
		{
			email:        "fakeuser_1@email.com",
			password:     "abc123password!",
			chirpReqBody: `{"body": "Hello, world"}`,
			expectedChirp: Chirp{
				Body: "Hello, world",
			},
		},
		{
			email:        "fakeuser_2@email.com",
			password:     "abc123password!",
			chirpReqBody: `{"body": "Hello, world"}`,
			expectedChirp: Chirp{
				Body: "Hello, world",
			},
		},
	}

	cfg := initApiConfig()
	for _, c := range cases {
		// Create new user
		_, _, err := createTestUser(cfg, c.email, c.password)
		if err != nil {
			t.Error(err)
		}

		// Log in to get the auth token required for posting
		loginUserBody := fmt.Sprintf(`{"email": "%v","password": "%v"}`, c.email, c.password)
		request := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginUserBody))
		w := httptest.NewRecorder()

		cfg.handlerLogin()(w, request)

		loggedInUser := &LoginResponse{}
		decoder := json.NewDecoder(w.Result().Body)
		err = decoder.Decode(&loggedInUser)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		// Pass the auth token, create the chirp
		chirpRequest := httptest.NewRequest("POST", "/api/chirp", strings.NewReader(c.chirpReqBody))
		chirpRequest.Header.Add("Authorization", "Bearer "+loggedInUser.Token)
		w = httptest.NewRecorder()

		cfg.postChirpHandler()(w, chirpRequest)

		// Get chirp response. validate fields
		chirpResp := &Chirp{}
		decoder = json.NewDecoder(w.Result().Body)
		err = decoder.Decode(&chirpResp)
		if err != nil {
			t.Error(err)
		}

		// Verify knowm and dynamic fields
		if chirpResp.Body != c.expectedChirp.Body {
			t.Error(formatTestError(*chirpResp, chirpResp.Body, c.expectedChirp.Body))
		}

		if err := uuid.Validate(string(chirpResp.ID.String())); err != nil {
			t.Error(formatTestError(*chirpResp, err, "a UUID"))
		}

		if chirpResp.UserID != loggedInUser.ID {
			t.Error(formatTestError(*chirpResp, chirpResp.UserID, loggedInUser.ID))
		}

		if chirpResp.CreatedAt.IsZero() {
			t.Error(formatTestError(*chirpResp, chirpResp.CreatedAt.IsZero(), "createdAt set to real timestamp"))
		}

		if chirpResp.UpdatedAt.IsZero() {
			t.Error(formatTestError(*chirpResp, chirpResp.UpdatedAt.IsZero(), "updatedAt set to real timestamp"))
		}
	}
}

func TestGetChirps(t *testing.T) {
	cfg := initApiConfig()

	// Create users
	users, passwords, err := createMultipleUsers(cfg, 3)
	if err != nil {
		t.Error(err)
	}

	// Create message corresponding to each user
	messages := []string{}
	for i := range len(users) {
		messages = append(messages, fmt.Sprintf("Hello from user %v", i))
	}

	// Login with each user, post a chirp
	for i, u := range users {
		loginResp, err := loginUser(cfg, u.Email, passwords[i])
		if err != nil {
			t.Error(err)
		}

		postChirp(cfg, loginResp.Token, messages[i])
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
	if len(chirps) != len(messages) {
		t.Errorf("failed to get all chirps: %v", chirps)
	}

	// All messages listed
	for _, message := range messages {
		if !slices.ContainsFunc(chirps, func(c Chirp) bool {
			return c.Body == message
		}) {
			t.Errorf("missing messages %v, Chirps: %v", messages, chirps)
		}
	}
}

func TestDeleteChirp(t *testing.T) {
	cfg := initApiConfig()

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
