package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	cfg := initApiConfig()

	// Clear users, create new user
	deleteAllUsers(cfg, t)

	email := "fakeuser@email.com"
	password := "abc123password!"
	createUserBody := fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, password)

	user, _, err := createTestUser(cfg, createUserBody)
	if err != nil {
		t.Errorf("failed to create user %v: %v", createUserBody, err)
		t.FailNow()
	}

	// Login
	loginUserBody := fmt.Sprintf(`{"email": "%v","password": "%v"}`, email, password)
	request := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginUserBody))
	w := httptest.NewRecorder()

	cfg.handlerLogin()(w, request)

	// Read response into struct
	loggedInUser := &LoginResponse{}
	decoder := json.NewDecoder(w.Result().Body)
	err = decoder.Decode(&loggedInUser)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify fields
	if loggedInUser.Email != user.Email {
		t.Errorf("expected: %v\nactual: %v", loggedInUser.Email, user.Email)
	}

	if loggedInUser.Token == "" {
		t.Errorf("missing jwt token:\n\t%v", loggedInUser)
	}
}
