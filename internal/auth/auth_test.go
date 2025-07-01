package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	// 1. Hash not the same as original password
	input := "abc123password"
	actual, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	assertNotEqual(actual, input, input, t)

	// 2. Hash is not empty
	input = "abc123password"
	assertNotEqual(actual, "", input, t)

	// 3. Hash is salted (different hash for same passwords)
	input = "abc123"
	actual1, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	actual2, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	assertNotEqual(actual1, actual2, input, t)
}

func TestCheckPassword(t *testing.T) {
	// Verify plaintext password matches hashed password
	plaintextPassword := "abc123password"
	hashedPassword, err := HashPassword(plaintextPassword)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEqual(CheckPasswordHash(plaintextPassword, hashedPassword), nil, nil, t)
}

func TestJWTGeneratesAToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration("1h")

	token, err := MakeJWT(userID, tokenSecret, expiresIn)

	if err != nil {
		t.Errorf("%v", err)
	}

	assertNotEqual(token, "", nil, t)
}

func TestValidateTokenExtractsUserID(t *testing.T) {
	// Create token, embedding userID
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration("1h")

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Extract userID from token
	extractedUserID, err := ValidateToken(token, tokenSecret)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify it's the same token
	assertEqual(extractedUserID, userID, nil, t)
}

func TestValidateTokenRejectExpiredToken(t *testing.T) {
	// Create new token with extremly short duration
	userID := uuid.New()
	tokenSecret := "acb123xyz!@#"
	expiresIn, _ := time.ParseDuration("1ns")

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Allow token to expire
	// TODO Replace sleep
	expirationDuration, _ := time.ParseDuration("5ns")
	time.Sleep(expirationDuration)

	_, err = ValidateToken(token, tokenSecret)

	// Check err message for invalid token and invalid claims error messages
	// https://github.com/golang-jwt/jwt/blob/v5.2.2/errors.go#L8
	expectedErrors := []error{
		jwt.ErrTokenExpired,
		jwt.ErrTokenInvalidClaims,
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
	expiresIn, _ := time.ParseDuration("1h")

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
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

func assertEqual(first, second, input any, t *testing.T) {
	if first != second {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, first, second)
	}
}

func assertNotEqual(first, second, input any, t *testing.T) {
	if first == second {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t*not* %v", input, first, second)
	}
}
