package api

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HealthHandler)
	mux.HandleFunc("GET /timezones", TimezonesHandler)
	mux.HandleFunc("GET /version", VersionHandler)

	loggingMux := RequestLogger(mux)

	return loggingMux
}