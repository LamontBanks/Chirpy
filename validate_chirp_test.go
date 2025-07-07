package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCensoredBannedWords(t *testing.T) {
	input := "What the sharbert? This is fornax crazy. Really, it's FORNAXing crazy kerfuffle"
	expected := "What the sharbert? This is **** crazy. Really, it's FORNAXing crazy ****"

	actual := censoredBannedWords(input)

	assertEquals(actual, expected, input, t)
}

func TestAllBannedWords(t *testing.T) {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	expected := "****"

	for _, word := range bannedWords {
		censoredWord := censoredBannedWords(word)

		assertEquals(censoredWord, expected, word, t)
	}
}

func TestValidateChirpHandler(t *testing.T) {
	// Chirp posts normally
	inputBody := `{"body":"I had something interesting for breakfast"}`
	expectedBody := `{"body":"I had something interesting for breakfast"}`
	expectedStatusCode := http.StatusOK

	// Request
	request := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(inputBody))
	request.Header.Set("Content-Type", "application/json")
	httpRecorder := httptest.NewRecorder()

	validateChirpHandler(httpRecorder, request)

	// Response
	response := httpRecorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Assertions
	assertEquals(response.StatusCode, expectedStatusCode, nil, t)
	assertEquals(string(body), expectedBody, nil, t)
}

func TestTooLongChirp(t *testing.T) {
	inputBody := `{"body": "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."}`
	expectedStatusCode := http.StatusBadRequest
	expectedBody := `{"error":"Chirp is too long"}`

	request := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(inputBody))
	request.Header.Set("Content-Type", "application/json")

	httpRecorder := httptest.NewRecorder()
	validateChirpHandler(httpRecorder, request)

	response := httpRecorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEquals(response.StatusCode, expectedStatusCode, nil, t)
	assertEquals(string(body), expectedBody, inputBody, t)
}

// Check if actual == expected
// `input` is an optional includsion of the original input
func assertEquals(actual, expected, input any, t *testing.T) {
	if actual != expected {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, actual, expected)
	}
}
