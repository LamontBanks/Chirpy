package handlers

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCensoredBannedWords(t *testing.T) {
	input := "What the sharbert? This is fornax crazy. Really, it's FORNAXing crazy kerfuffle"
	expected := "What the sharbert? This is **** crazy. Really, it's FORNAXing crazy ****"

	actual := censoredBannedWords(input)

	assertEqual(expected, actual, input, t)
}

func TestValidateChirpHandler(t *testing.T) {
	input := `{"body":"test kerfuffle"}`
	expected := `{"body":"test ****"}`

	// Request
	request := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(input))
	request.Header.Set("Content-Type", "application/json")

	// Call handler
	httpRecorder := httptest.NewRecorder()
	ValidateChirpHandler(httpRecorder, request)

	// Read response
	response := httpRecorder.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Assert status code, response body, etc.
	assertEqual(string(body), expected, input, t)
}

func assertEqual(actual, expected, input any, t *testing.T) {
	if actual != expected {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, actual, expected)
	}
}
