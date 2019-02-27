package http

import (
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
)

func GetDeck(decks store.DeckStore) StoreGetFunc {
	return func(id uuid.UUID) (interface{}, error) { return decks.Get(id) }
}
