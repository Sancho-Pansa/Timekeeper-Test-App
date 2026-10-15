package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
	"github.com/Sancho-Pansa/Timekeeper/internal/store"
)

func TestValidateOffset(t *testing.T) {
	tests := []struct{
		name string
		input string
		want int
		wantError bool
	} {
		{ "UTC", "0", 0, false },
		{ "Positive offset", "7", 7, false },
		{ "Negative offset", "-4", -4, false },
		{ "Minimum limit", "-12", -12, false },
		{ "Maximum limit", "12", 12, false },
		{ "Out of bounds 01", "128", 0, true },
		{ "Out of bounds 02", "-128", 0, true },
		{ "Invalid input", "The Word", 0, true },
		{ "Empty string", "", 0, true },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateOffset(tt.input)
			if (err != nil) != tt.wantError {
				t.Fatalf("Error in 'validateOffset(%s) = %v; wantError = %v'", tt.input, err, tt.wantError)
			}

			if (err == nil) && got != tt.want {
				t.Errorf("Error in validateOffset(%s). Got: %v, Expected: %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	authService := auth.NewService("user", "password")
	memoryStore := store.NewMemoryStore()

	router := NewRouter(authService, memoryStore)

	request := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Endpoint /health :: Status: %d, want %d", response.StatusCode, http.StatusOK)
	}

	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("Endpoint /health :: Failed to decode response: %v", err)
	}

	if body.Status != "ok" {
		t.Errorf("Endpoint /health :: Status = %q, want %q", body.Status, "ok")
	}
}

func TestGetWithoutToken(t *testing.T) {
	authService := auth.NewService("user", "password")
	memoryStore := store.NewMemoryStore()
	router := NewRouter(authService, memoryStore)

	request := httptest.NewRequest(
		http.MethodGet,
		"/favorites",
		nil,
	)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Enpoint /favorites :: Status %d, expected %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticatedProfile(t *testing.T) {
	authService := auth.NewService("user", "password")
	memoryStore := store.NewMemoryStore()
	router := NewRouter(authService, memoryStore)

	loginBody := strings.NewReader(`{
		"username": "user",
		"password": "password"
	}`)

	loginRequest := httptest.NewRequest(
		http.MethodPost,
		"/auth",
		loginBody,
	)
	loginRequest.Header.Set("Content-Type", "application/json")
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginRequest)

	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("Endpoint /auth :: Login failed, Status %d, Body: %s", loginRecorder.Code, loginRecorder.Body)
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(loginRecorder.Body.Bytes(), &authResponse); err != nil {
		t.Fatalf("Endpoint /auth :: Failed to decode token: %v", err)
	}

	if authResponse.AccessToken == "" {
		t.Fatal("Endpoint /auth :: Empty token")
	}

	profileRequest := httptest.NewRequest(http.MethodGet, "/profile", nil)
	profileRequest.Header.Set("Authorization", "Bearer " + authResponse.AccessToken)
	profileRecorder := httptest.NewRecorder()
	router.ServeHTTP(profileRecorder, profileRequest)

	if profileRecorder.Code != http.StatusOK {
		t.Fatalf("Endpoint /profile :: Reqiest failed, HTTP %d: %s", profileRecorder.Code, profileRecorder.Body.String())
	}

	var profileResponse ProfileResponse
	if err := json.Unmarshal(profileRecorder.Body.Bytes(), &profileResponse); err != nil {
		t.Fatalf("Endpoint /profile :: JSON decode error, %v", err)
	}

	if profileResponse.Username != "user" {
		t.Errorf("Endpoint /profile :: Got Username %s, expected 'user'", profileResponse.Username)
	}
}