package http

import (
	"encoding/json"
	"errors"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

func DecodeCard(r *http.Request) (interface{}, error) {
	u := &store.Card{}

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		return nil, ErrCantParseJson
	}

	return u, nil
}

func PostCard(cards store.CardStore) StorePostFunc {
	return func(data interface{}) error {
		if card, ok := data.(*store.Card); !ok {
			return errors.New("not a Card")
		} else {
			return cards.Add(card)
		}
	}
}

func GetCard(cards store.CardStore) StoreGetFunc {
	return func(id string) (interface{}, error) { return cards.Get(id) }
}

func PutCard(cards store.CardStore) StorePutFunc {
	return func(data interface{}) error {
		if card, ok := data.(*store.Card); !ok {
			return errors.New("not a Card")
		} else {
			return cards.Save(card)
		}
	}
}

func DeleteCard(cards store.CardStore) StoreDeleteFunc {
	return func(id string) error { return cards.Delete(id) }
}
