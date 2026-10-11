package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/api"
	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
)

func main() {
	username := os.Getenv("APP_USERNAME")
	password := os.Getenv("APP_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("APP_USERNAME and APP_PASSWORD environmet vars are required")
	}

	authService := auth.NewService(username, password)

	server := &http.Server {
		Addr: ":8080",
		Handler: api.NewRouter(authService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Timekeeper is listening on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}