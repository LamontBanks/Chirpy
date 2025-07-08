package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCensoredBannedWords(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single word censored",
			input:    "kerfuffle",
			expected: "****",
		},
		{
			name:     "Single word censored",
			input:    "sharbert",
			expected: "****",
		},
		{
			name:     "Single word censored",
			input:    "fornax",
			expected: "****",
		},
		{
			name:     "Mixed input",
			input:    "What the sharbert? This is fornax crazy. Really, it's FORNAXing crazy kerfuffle",
			expected: "What the sharbert? This is **** crazy. Really, it's FORNAXing crazy ****",
		},
	}

	for _, c := range cases {
		actual := censoredBannedWords(c.input)
		if actual != c.expected {
			t.Error(formatTestError(c.name, c.expected, actual))
		}
	}
}

func TestValidateChirpHandler(t *testing.T) {
	cases := []struct {
		name               string
		inputBody          string
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name:               "No censoring",
			inputBody:          `{"body":"I had something interesting for breakfast"}`,
			expectedBody:       `{"body":"I had something interesting for breakfast"}`,
			expectedStatusCode: http.StatusOK,
		},
		{

			name:               "Censor words, ignore unexpected input",
			inputBody:          `{"body":"I hear Mastodon is better than Chirpy. sharbert I need to migrate","extra": "this should be ignored"}`,
			expectedBody:       `{"body":"I hear Mastodon is better than Chirpy. **** I need to migrate"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Censor word in string",
			inputBody:          `{"body":"I really need a kerfuffle to go to bed sooner, Fornax !"}`,
			expectedBody:       `{"body":"I really need a **** to go to bed sooner, **** !"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Input too long",
			inputBody:          `{"body":"lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."}`,
			expectedBody:       `{"error":"Chirp is too long"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		// Request
		request := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(c.inputBody))
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
		if response.StatusCode != c.expectedStatusCode {
			t.Error(formatTestError(c.name, c.expectedBody, response.StatusCode))
		}

		if string(body) != c.expectedBody {
			t.Error(formatTestError(c.name, c.expectedBody, string(body)))
		}
	}
}

// Check if actual == expected
// `input` is an optional includsion of the original input
func assertEquals(actual, expected, input any, t *testing.T) {
	if actual != expected {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, actual, expected)
	}
}

func formatTestError(testname, actual, expected any) string {
	return fmt.Sprintf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", testname, actual, expected)
}
