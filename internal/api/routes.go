package api

import (
	"net/http"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
	"github.com/Sancho-Pansa/Timekeeper/internal/store"
)

func NewRouter(authService *auth.Service, favouriteStore *store.MemoryStore) http.Handler {
	mux := http.NewServeMux()

	h := &Handler{
		Auth:  authService,
		Store: favouriteStore,
	}

	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /timezones", TimezonesHandler)
	mux.HandleFunc("GET /version", VersionHandler)
	mux.HandleFunc("POST /auth", h.AuthHandler)

	mux.Handle("GET /profile", RequireAuthorization(authService, http.HandlerFunc(h.ProfileHandler)))
	mux.Handle("GET /favorites", RequireAuthorization(authService, http.HandlerFunc(h.GetFavorites)))
	mux.Handle("POST /favorites/{offset}", RequireAuthorization(authService, http.HandlerFunc(h.CreateFavorite)))
	mux.Handle("DELETE /favorites/{offset}", RequireAuthorization(authService, http.HandlerFunc(h.DeleteFavorites)))

	loggingMux := RequestLogger(mux)

	return loggingMux
}
