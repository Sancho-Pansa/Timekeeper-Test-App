package api

import(
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now()
		log.Printf("Request received: %s :: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Handled in %v", time.Since(timestamp))
	})
}