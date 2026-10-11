package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
	"github.com/Sancho-Pansa/Timekeeper/internal/timezones"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType string `json:"tokenType"`
	ExpiresIn int `json:"expiresIn"`
}

type Handler struct {
	Auth *auth.Service
}

func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()

	var request AuthRequest

	if err := jsonDecoder.Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var extra any
	if err := jsonDecoder.Decode(&extra); err != io.EOF {
		http.Error(w, "Expected only one JSON object", http.StatusBadRequest)
		return
	}

	if !h.Auth.Authenticate(request.Username, request.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.Auth.CreateSession(request.Username)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJson(w, http.StatusOK, AuthResponse {
		AccessToken: token,
		TokenType: "Bearer",
		ExpiresIn: int(auth.TokenTTL.Seconds()),
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writeJson(w, http.StatusOK, map[string]string {
		"username": username,
	})
}

func TimezonesHandler(w http.ResponseWriter, r *http.Request) {
	result := timezones.GenerateTimezoneList()
	writeJson(w, http.StatusOK, result)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]string{
		"version": "0.0.5",
		"time": time.Now().Format(time.RFC3339),
	})
}