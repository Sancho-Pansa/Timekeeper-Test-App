package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/api"
)

func main() {
	server := &http.Server {
		Addr: ":8080",
		Handler: api.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Timekeeper is listening on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}