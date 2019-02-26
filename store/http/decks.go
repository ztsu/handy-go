package http

import "github.com/ztsu/handy-go/store"

func GetDeck(decks store.DeckStore) StoreGetFunc {
	return func(id store.UUID) (interface{}, error) { return decks.Get(id) }
}
