package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
)

// Отдельный тип, чтобы избежать коллизии ключей с дргуими пакетами
type contextKey string

const usernameKey contextKey = "username"

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		log.Printf("Request received: %s :: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Handled in %v", time.Since(timestamp))
	})
}

func RequireAuthorization(authService *auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(w)
			return
		}

		username, valid := authService.ValidateToken(parts[1])
		if !valid {
			unauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), usernameKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}