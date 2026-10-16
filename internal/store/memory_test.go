package store

import (
	"reflect"
	"sync"
	"testing"
)

func TestGetNewFavorites(t *testing.T) {
	testStore := NewMemoryStore()
	username := "user"
	favoriteOffset := -4
	favorites := []int { -4 }

	t.Run("Newly created store is empty", func(t *testing.T) {
		userStore := testStore.GetFavourites(username)
		if len(userStore) != 0 {
			t.Error("New Store is not Empty!")
		}
	})

	t.Run("CRUD operations with offsets", func(t *testing.T) {
		if isAdded := testStore.AddFavourite(username, favoriteOffset); !isAdded {
			t.Error("Failed to add new offset entry.")
		}

		if isAdded := testStore.AddFavourite(username, favoriteOffset); isAdded {
			t.Error("Received 'true' on add existing entry; expected 'false'.")
		}

		if !reflect.DeepEqual(favorites, testStore.GetFavourites(username)) {
			t.Errorf("Got favorites: %v; Expected: %v", testStore.GetFavourites(username), favorites)
		}

		if !testStore.DeleteFavourite(username, favoriteOffset) {
			t.Error("Delete favourite offset failed. Got 'false' from operation, expected 'true'")
		}

		if testStore.DeleteFavourite(username, favoriteOffset) {
			t.Error("Error in repeated favourite offset deletion. Got 'true' from operation, expected 'false'")
		}
	})
}

func TestMemoryStoreConcurrency(t *testing.T) {
	s := NewMemoryStore()
	username := "user"

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			offset := i % 25 - 12
			s.AddFavourite(username, offset)
			s.GetFavourites(username)
			s.DeleteFavourite(username, offset)
		}(i)
	}

	wg.Wait()
}