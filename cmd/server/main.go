package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sancho-Pansa/Timekeeper/internal/api"
	"github.com/Sancho-Pansa/Timekeeper/internal/auth"
	"github.com/Sancho-Pansa/Timekeeper/internal/store"
)

func main() {
	username := os.Getenv("APP_USERNAME")
	password := os.Getenv("APP_PASSWORD")
	port := os.Getenv("APP_PORT")

	if username == "" || password == "" {
		log.Fatal("APP_USERNAME and APP_PASSWORD environmet vars are required")
	}
	if port == "" {
		port = "8080"
	}

	authService := auth.NewService(username, password)
	memoryStore := store.NewMemoryStore()

	server := &http.Server {
		Addr: ":" + port,
		Handler: api.NewRouter(authService, memoryStore),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Настройка реакции на получение SIGTERM от ОС ("подписка")
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	/* Если оставить ListenAndServe() в обычном потоке, то дальше него код не пойдет, в
		том числе и ожидание сигнала останова.
		А свои ошибки он передает в канал serverErrors
	 */
	go func() {
		log.Printf("Timekeeper is listening on :%s", port)
		serverErrors <- server.ListenAndServe()
	}()

	select { // как switch, но для обработки действий с каналами
	case <- ctx.Done(): // Т.е. контекст отменен
		log.Println("Shutdown signal received")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server failed: %v", err)
		}
		return
	}
	// Пока что-то из select не произойдет, код дальше не пойдет

	stop()

	// Ещё один контекст с дедлайном в 10 секунд - на случай, если сервер занят очень крупным запросом
	// Через 10 секунд он автоматически завершится (в отличие от контекста выше, который среагирует только на SIGTERM)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdwon has failed: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("Force close has failed: %v", err)
		}
		os.Exit(1)
	}

	log.Println("HTTP server successfully stopped")
}