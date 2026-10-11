package api

import (
	"net/http"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
)

func NewRouter(authService *auth.Service) http.Handler {
	mux := http.NewServeMux()

	h := &Handler{
		Auth: authService,
	}

	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /timezones", TimezonesHandler)
	mux.HandleFunc("GET /version", VersionHandler)
	mux.HandleFunc("POST /auth", h.AuthHandler)

	loggingMux := RequestLogger(mux)

	return loggingMux
}