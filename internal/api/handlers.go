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

type ProfileResponse struct {
	Username         string    `json:"username"`
	FavouriteOffsets []int     `json:"favoriteOffsets"`
	ServerTime       time.Time `json:"serverTime"`
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

	if request.Username == "" || request.Password == "" {
		http.Error(w, "Empty credentials", http.StatusUnauthorized)
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

	numericOffset, err := validateOffset(pathOffset)
	if err != nil {
		writeValidationError(w, r, err)
	}

	favOffset := FavouriteOffset{Offset: &numericOffset}

	username, ok := r.Context().Value(usernameKey).(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdded := h.Store.AddFavourite(username, *favOffset.Offset)
	if !isAdded {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	favorites := h.Store.GetFavourites(username)

	writeJson(w, http.StatusOK, favorites)
}

func (h *Handler) DeleteFavorites(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pathOffset := r.PathValue("offset")
	numericOffset, err := validateOffset(pathOffset)
	if err != nil {
		writeValidationError(w, r, err)
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

	writeJson(w, http.StatusOK, ProfileResponse{
		Username:         username,
		FavouriteOffsets: h.Store.GetFavourites(username),
		ServerTime:       time.Now().UTC(),
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

func validateOffset(offsetValue string) (int, error) {
	numericOffset, err := strconv.Atoi(offsetValue)
	if err != nil {
		return 0, fmt.Errorf("Numeric offset is required")
	}
	if numericOffset < -12 || numericOffset > 12 {
		return 0, fmt.Errorf("Offset must be between -12 and 12")
	}
	return numericOffset, nil
}

func writeValidationError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, fmt.Sprintf("Error in '%s' method: %s", r.Pattern, err.Error()), http.StatusBadRequest)
}
