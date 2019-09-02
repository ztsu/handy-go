package http

import (
	"encoding/json"
	"errors"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

func DecodeDeck(r *http.Request) (interface{}, error) {
	u := &store.Deck{}

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		return nil, ErrCantParseJson
	}

	return u, nil
}

func PostDeck(decks store.DeckStore) StorePostFunc {
	return func(data interface{}) error {
		if deck, ok := data.(*store.Deck); !ok {
			return errors.New("not a Deck")
		} else {
			return decks.Add(deck)
		}
	}
}

func GetDeck(decks store.DeckStore) StoreGetFunc {
	return func(id string) (interface{}, error) { return decks.Get(id) }
}

func PutDeck(decks store.DeckStore) StorePutFunc {
	return func(data interface{}) error {
		if deck, ok := data.(*store.Deck); !ok {
			return errors.New("not a Deck")
		} else {
			return decks.Save(deck)
		}
	}
}

func DeleteDeck(decks store.DeckStore) StoreDeleteFunc {
	return func(id string) error { return decks.Delete(id) }
}
