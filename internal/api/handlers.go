package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
	"github.com/Sancho-Pansa/Timekeeper/internal/store"
	"github.com/Sancho-Pansa/Timekeeper/internal/timezones"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

type Handler struct {
	Auth  *auth.Service
	Store *store.MemoryStore
}

type FavouriteOffset struct {
	Offset *int `json:"offset"`
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

	writeJson(w, http.StatusOK, AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(auth.TokenTTL.Seconds()),
	})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func (h *Handler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	pathOffset := r.PathValue("offset")

	numericOffset, err := strconv.Atoi(pathOffset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in '%s' method: numeric offset is required", r.Pattern), http.StatusBadRequest)
		return
	}

	if numericOffset < -12 || numericOffset > 12 {
		http.Error(w, fmt.Sprintf("Error in '%s' method: Offset must be between -12 and 12", r.Pattern), http.StatusBadRequest)
		return
	}

	favOffset := FavouriteOffset{Offset: &numericOffset}

	username, ok := r.Context().Value(usernameKey).(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdded := h.Store.AddFavourite(username, *favOffset.Offset)
	if !isAdded {
		http.Error(w, "This offset is already in Favourites", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)

	favorites := h.Store.GetFavourites(username)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writeJson(w, http.StatusOK, favorites)
}

func (h *Handler) DeleteFavorites(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathOffset := r.PathValue("offset")
	numericOffset, ok := validateOffset(pathOffset, w, r)
	if !ok {
		return
	}

	isDeleted := h.Store.DeleteFavourite(username, numericOffset)
	if !isDeleted {
		http.Error(w, fmt.Sprintf("Error in '%s' method: Offset not found", r.Pattern), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{
		"username": username,
	})
}

func TimezonesHandler(w http.ResponseWriter, r *http.Request) {
	result := timezones.GenerateTimezoneList()
	writeJson(w, http.StatusOK, result)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]string{
		"version": os.Getenv("APP_VERSION"),
		"time":    time.Now().Format(time.RFC3339),
	})
}

func validateOffset(offsetValue string, w http.ResponseWriter, r *http.Request) (int, bool) {
	numericOffset, err := strconv.Atoi(offsetValue)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in '%s' method: numeric offset is required", r.Pattern), http.StatusBadRequest)
		return 0, false
	}
	if numericOffset < -12 || numericOffset > 12 {
		http.Error(w, fmt.Sprintf("Error in '%s' method: Offset must be between -12 and 12", r.Pattern), http.StatusBadRequest)
		return 0, false
	}
	return numericOffset, true
}