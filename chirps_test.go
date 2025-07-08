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

	// Create multiple users
	users, passwords, err := createTestUsers(cfg, 3)
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

		_, err = postChirp(cfg, loginResp.Token, messages[i])
		if err != nil {
			t.Error(err)
		}
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
		t.Errorf("failed to get all chirps: %+v", chirps)
	}

	// All messages listed
	for _, msg := range messages {
		if !slices.ContainsFunc(chirps, func(c Chirp) bool {
			return c.Body == msg
		}) {
			t.Errorf("missing message '%v': %+v", msg, chirps)
		}
	}
}

func TestDeleteChirp(t *testing.T) {
	cases := []struct {
		name                    string
		messages                []string
		chirpsToDelete          []Chirp
		expectedRemainingChirps []Chirp
		expectedRespCode        int
	}{
		{
			name:     "Delete single chirp",
			messages: []string{"abc", "123", "xyz"},
			chirpsToDelete: []Chirp{
				{
					Body: "abc",
				},
			},
			expectedRespCode: http.StatusNoContent,
			expectedRemainingChirps: []Chirp{
				{
					Body: "123",
				},
				{
					Body: "xyz",
				},
			},
		},
	}

	cfg := initApiConfig()

	for _, c := range cases {
		err := deleteAllUsersAndPosts(cfg)
		if err != nil {
			t.Error(err)
		}

		// Create single user
		users, passwords, err := createTestUsers(cfg, 1)
		if err != nil {
			t.Error(err)
		}

		// Login, post chirps
		loggedInUser, err := loginUser(cfg, users[0].Email, passwords[0])
		postedChirps := []Chirp{}

		for _, msg := range c.messages {
			if err != nil {
				t.Error(err)
			}

			chirp, err := postChirp(cfg, loggedInUser.Token, msg)
			if err != nil {
				t.Error(err)
			}

			postedChirps = append(postedChirps, chirp)
		}

		// Delete specified chirps
		for _, chirpToDelete := range c.chirpsToDelete {
			// Find the chirp based on the message text
			index := slices.IndexFunc(postedChirps, func(chirp Chirp) bool {
				return chirp.Body == chirpToDelete.Body
			})

			deleteChirpRequest := httptest.NewRequest("DELETE", "/api/chirps/", nil)
			deleteChirpRequest.SetPathValue("chirpID", postedChirps[index].ID.String())
			deleteChirpRequest.Header.Add("Authorization", "Bearer "+loggedInUser.Token)

			w := httptest.NewRecorder()

			cfg.deleteChirpHandler()(w, deleteChirpRequest)

			// DELETE Success response
			if w.Result().StatusCode != c.expectedRespCode {
				t.Error(formatTestError(w, w.Result().StatusCode, c.expectedRespCode))
			}

			// Ensure chirp was removed
			getChirpsReq := httptest.NewRequest("GET", "/api/chirps", nil)
			w = httptest.NewRecorder()
			cfg.getChirps()(w, getChirpsReq)

			actualRemainingChirps := []Chirp{}
			decoder := json.NewDecoder(w.Result().Body)
			err = decoder.Decode(&actualRemainingChirps)
			if err != nil {
				t.Error(err)
			}

			// Correct chirps (based on text) remain
			correctChirpsRemain := slices.EqualFunc(actualRemainingChirps, c.expectedRemainingChirps, func(actual, expected Chirp) bool {
				return actual.Body == expected.Body
			})
			if !correctChirpsRemain {
				t.Error(formatTestError("incorrect chirps remaining after DELETE", actualRemainingChirps, c.expectedRemainingChirps))
			}
		}
	}
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
