package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	HealthHandler(recorder, request)

	assertEqual(recorder.Result().StatusCode, http.StatusOK, request, t)
}
