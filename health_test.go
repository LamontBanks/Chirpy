package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name                 string
		method               string
		expectedResponseCode int
	}{
		{
			name:                 "GET",
			method:               http.MethodGet,
			expectedResponseCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		request := httptest.NewRequest(c.method, "/healthz", nil)
		recorder := httptest.NewRecorder()

		healthHandler(recorder, request)

		if recorder.Result().StatusCode != c.expectedResponseCode {
			t.Error(formatTestError(c.name, recorder.Result().StatusCode, c.expectedResponseCode))
		}
	}

}
