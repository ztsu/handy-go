package http

import (
	"encoding/json"
	"errors"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

func DecodeDeckCard(r *http.Request) (interface{}, error) {
	u := &store.DeckCard{}

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		return nil, ErrCantParseJson
	}

	return u, nil
}

func PostDeckCard(cards store.DeckCardStore) StorePostFunc {
	return func(data interface{}) error {
		if card, ok := data.(*store.DeckCard); !ok {
			return errors.New("not a DeckCard")
		} else {
			return cards.Add(card)
		}
	}
}

func GetDeckCard(cards store.DeckCardStore) StoreGetFunc {
	return func(id string) (interface{}, error) { return cards.Get(id) }
}

func PutDeckCard(cards store.DeckCardStore) StorePutFunc {
	return func(data interface{}) error {
		if card, ok := data.(*store.DeckCard); !ok {
			return errors.New("not a DeckCard")
		} else {
			return cards.Save(card)
		}
	}
}

func DeleteDeckCard(cards store.DeckCardStore) StoreDeleteFunc {
	return func(id string) error { return cards.Delete(id) }
}
