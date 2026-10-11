package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"sync"
	"time"
)

const TokenTTL = time.Hour

type Session struct {
	Username string
	ExpiresAt time.Time
}

type Service struct {
	username string
	password string

	mu sync.RWMutex
	sessions map[string]Session
}

func NewService(username, password string) *Service {
	return &Service{
		username: username,
		password: password,
		sessions: make(map[string]Session),
	}
}

func (s *Service) Authenticate(username, password string) bool {
	userHash := sha256.Sum256([]byte(username))
	expectedUserHash := sha256.Sum256([]byte(s.username))

	passwordHash := sha256.Sum256([]byte(password))
	expectedPasswordHash := sha256.Sum256([]byte(s.password))

	userOk := subtle.ConstantTimeCompare(userHash[:], expectedUserHash[:])
	passwordOk := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:])
	log.Println("Auth in progress")
	return userOk == 1 && passwordOk == 1
}

func (s *Service) CreateSession(username string) (string, error) {
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[token] = Session{
		Username: username,
		ExpiresAt: time.Now().Add(TokenTTL),
	}
	log.Println("Token created")

	return token, nil
}

func (s *Service) ValidateToken(token string) (string, bool) {
	s.mu.RLock()
	sessions, exists := s.sessions[token]
	s.mu.RUnlock()

	if !exists || !time.Now().Before(sessions.ExpiresAt) {
		return "", false
	}

	return sessions.Username, true
}