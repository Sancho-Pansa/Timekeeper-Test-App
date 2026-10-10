package api

import (
	"net/http"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/timezones"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func TimezonesHandler(w http.ResponseWriter, r *http.Request) {
	result := timezones.GenerateTimezoneList()
	writeJson(w, http.StatusOK, result)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, map[string]string{
		"version": "0.1.0",
		"time": time.Now().Format(time.RFC3339),
	})
}