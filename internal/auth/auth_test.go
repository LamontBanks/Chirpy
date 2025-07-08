package auth

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		notExpected string
	}{
		{
			name:        "Hash different from original password",
			input:       "abc123password",
			notExpected: "abc123password",
		},
		{
			name:        "Hash is not empty",
			input:       "abc123password",
			notExpected: "",
		},
	}

	for _, c := range cases {
		actual, err := HashPassword(c.input)
		if err != nil {
			t.Error(err)
		}

		if actual == c.notExpected {
			t.Error(formatTestError(c.name, actual, "not: "+c.notExpected))
		}
	}
}

func TestHashPasswordIsSalted(t *testing.T) {
	cases := []struct {
		input string
	}{
		{
			input: "abc123",
		},
	}

	for _, c := range cases {
		hash1, err := HashPassword(c.input)
		if err != nil {
			t.Error(err)
		}

		hash2, err := HashPassword(c.input)
		if err != nil {
			t.Error(err)
		}

		if hash1 == hash2 {
			t.Error(formatTestError("Hashed password different for same input",
				fmt.Errorf("%v = %v", hash1, hash2),
				fmt.Errorf("%v != %v", hash1, hash2)))
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		input string
	}{
		{
			input: "abc123password",
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Error(err)
		}

		actual := CheckPasswordHash(c.input, hash)

		if actual != nil {
			t.Error(formatTestError("No error returned", actual, nil))
		}
	}
}

func TestJWTGeneratesAToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration(JWT_TOKEN_DURATION)

	token, err := MakeJWT(userID, tokenSecret, expiresIn)

	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error(formatTestError(userID, token, `not ""`))
	}
}

func TestValidateTokenExtractsUserID(t *testing.T) {
	// Create token, embedding userID
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration(JWT_TOKEN_DURATION)

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	// Extract userID from token
	extractedUserID, err := ValidateToken(token, tokenSecret)
	if err != nil {
		t.Error(err)
	}

	// Verify it's the same token
	assertEqual(extractedUserID, userID, nil, t)
}

func TestValidateTokenRejectExpiredToken(t *testing.T) {
	// Create new token with extremly short duration
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, err := time.ParseDuration("1ns")
	if err != nil {
		t.Error(err)
	}

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	// Allow token to expire
	// TODO Replace sleep
	expirationDuration, _ := time.ParseDuration("5ns")
	time.Sleep(expirationDuration)

	userID, err = ValidateToken(token, tokenSecret)

	// Check err message for invalid token and invalid claims error messages
	// https://github.com/golang-jwt/jwt/blob/v5.2.2/errors.go#L8
	expectedErrors := []error{
		jwt.ErrTokenExpired,
		jwt.ErrTokenInvalidClaims,
	}

	if userID != uuid.Nil {
		t.Errorf("expired token returned userId %v, expected %v", userID, uuid.Nil)
	}

	for i := range expectedErrors {
		if !strings.Contains(err.Error(), expectedErrors[i].Error()) {
			t.Errorf("missing expected expired token Error: '%v'", expectedErrors[i])
		}
	}
}

func TestValidateTokenRejectWrongTokenSecret(t *testing.T) {
	// Create new token with
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration(JWT_TOKEN_DURATION)

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	// Attempt to read with wrong secret key
	wrongTokenSecret := "wrongSecret123"

	_, err = ValidateToken(token, wrongTokenSecret)

	// Check for expected error message(s)
	expectedErrors := []error{
		jwt.ErrTokenSignatureInvalid,
	}

	for i := range expectedErrors {
		if !strings.Contains(err.Error(), expectedErrors[i].Error()) {
			t.Errorf("missing expected invalid tokenSecret error: '%v'", expectedErrors[i])
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	// Extract correctly set bearer token
	expectedBearerToken := "abc123"

	header := httptest.NewRecorder().Header()
	header.Add("Authorization", "Bearer "+expectedBearerToken)

	actualToken, err := GetBearerToken(header)
	if err != nil {
		t.Error(err)
	}

	assertEqual(actualToken, expectedBearerToken, header, t)
}

func TestGetBearerTokenErrorIfNotSet(t *testing.T) {
	cases := []struct {
		name  string
		input []string
	}{
		{
			name:  "Missing 'Bearer ' prefix",
			input: []string{"Authorization", "abc123"},
		},
		{
			name:  "Misformatted Bearer prefix (no space)",
			input: []string{"Authorization", "Bearerabc123"},
		},
		{
			name:  "Missing token",
			input: []string{"Authorization", "Bearer"},
		},
	}

	for _, c := range cases {
		header := httptest.NewRecorder().Header()
		header.Add(c.input[0], c.input[1])

		_, actual := GetBearerToken(header)

		if actual == nil {
			t.Error(formatTestError(c.name, actual, fmt.Sprintf("%v: %v -> error", c.input[0], c.input[1])))
		}
	}
}

func TestGetAPIKey(t *testing.T) {
	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"Authorization", "ApiKey abc123"},
			expected: "abc123",
		},
		{
			input:    []string{"Authorization", "ApiKey a-b-c-1-2-3"},
			expected: "a-b-c-1-2-3",
		},
	}

	for _, c := range cases {
		header := httptest.NewRecorder().Header()
		header.Add(c.input[0], c.input[1])

		actualAPIKey, err := GetAPIKey(header)
		if err != nil {
			t.Error(err)
		}

		if actualAPIKey != c.expected {
			t.Error(formatTestError(c.input, actualAPIKey, c.expected))
		}
	}
}

func TestGetAPIKeyErrorIfNotSet(t *testing.T) {
	cases := []struct {
		name  string
		input []string
	}{
		{
			name:  "Missing 'ApiKey ' prefix",
			input: []string{"Authorization", "abc123"},
		},
		{
			name:  "Misformatted ApiKey prefix (no space)",
			input: []string{"Authorization", "ApiKeyabc123"},
		},
		{
			name:  "Missing token",
			input: []string{"Authorization", "ApiKey"},
		},
	}

	for _, c := range cases {
		header := httptest.NewRecorder().Header()
		header.Add(c.input[0], c.input[1])

		_, actual := GetAPIKey(header)

		if actual == nil {
			t.Error(formatTestError(c.name, actual, fmt.Sprintf("%v: %v -> error", c.input[0], c.input[1])))
		}
	}
}

func assertEqual(first, second, input any, t *testing.T) {
	if first != second {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, first, second)
	}
}

func formatTestError(testname, actual, expected any) string {
	return fmt.Sprintf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", testname, actual, expected)
}
