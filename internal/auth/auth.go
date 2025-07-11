package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Custom JWT Token with Claims
// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#NewWithClaims
// jwt.RegisterdClaims can used on its own, but wrapping in a custom Claims is the typical use case
// This allows additional public/private data to be added to the token
// Currently not adding any add'l info, but will follow the recommended usage:
// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-NewWithClaims-CustomClaimsType
type CustomClaims struct {
	jwt.RegisteredClaims
}

const (
	JWT_TOKEN_DURATION     = "1h"
	REFRESH_TOKEN_DURATION = "1440h"
)

// Returns `bcrypt`-hashed password
func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// Returns nil if plaintext password matches the hashed password, error otherwise
func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Returns a JSON Web Token (JWT) for the given user
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if userID == uuid.Nil || tokenSecret == "" {
		return "", fmt.Errorf("invalid userID: %v, tokenSecret: %v", userID, tokenSecret)
	}

	// Create token
	claims := CustomClaims{
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		},
	}

	// Create token, sign with given method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))

	return signedToken, err
}

// Validates the token, extracts and returns the userID
func ValidateToken(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

// Return the token from the `Authorization: Bearer <token>` header
func GetBearerToken(header http.Header) (string, error) {
	bearerToken := header.Get("Authorization")

	if bearerToken == "" {
		return "", fmt.Errorf("authorization header value not found")
	}

	// Included whitespace in split is intentional
	splitBearerToken := strings.Split(bearerToken, "Bearer ")

	if len(splitBearerToken) < 2 {
		return "", fmt.Errorf("invalid authorization header value")
	}

	bearerToken = strings.TrimSpace(splitBearerToken[1])

	return bearerToken, nil
}

// Return the key from the `Authorization: ApiKey <key>` header
func GetAPIKey(header http.Header) (string, error) {
	apiKey := header.Get("Authorization")

	if apiKey == "" {
		return "", fmt.Errorf("authorization header value not found")
	}

	// Included whitespace in split is intentional
	splitAPIKey := strings.Split(apiKey, "ApiKey ")

	if len(splitAPIKey) < 2 {
		return "", fmt.Errorf("invalid authorization header value")
	}

	apiKey = strings.TrimSpace(splitAPIKey[1])

	return apiKey, nil
}

// Returns a 256-bit, hex-encoded string
func MakeRefreshToken() (string, error) {
	randBits := make([]byte, 256)

	_, err := rand.Read(randBits)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(randBits)

	return token, nil
}
