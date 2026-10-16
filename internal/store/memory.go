package store

import (
	"sort"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex

	favourites map[string]map[int]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		favourites: make(map[string]map[int]struct{}),
	}
}

func (s *MemoryStore) GetFavourites(username string) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, 0)

	for offset := range s.favourites[username] {
		result = append(result, offset)
	}

	sort.Ints(result)

	return result
}

func (s *MemoryStore) AddFavourite(username string, offset int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.favourites[username] == nil {
		s.favourites[username] = make(map[int]struct{})
	}

	if _, exists := s.favourites[username][offset]; exists {
		return false
	}

	s.favourites[username][offset] = struct{}{}
	return true
}

func (s *MemoryStore) DeleteFavourite(username string, offset int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.favourites[username][offset]; !exists {
		return false
	}

	delete(s.favourites[username], offset)
	return true
}