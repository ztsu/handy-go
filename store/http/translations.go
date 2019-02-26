package http

import (
	"encoding/json"
	"errors"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

func GetTranslation(translations store.TranslationStore) StoreGetFunc {
	return func(id store.UUID) (interface{}, error) { return translations.Get(id) }
}

func DeleteTranslation(translations store.TranslationStore) StoreDeleteFunc {
	return func(id store.UUID) error { return translations.Delete(id) }
}

func DecodeTranslation(r *http.Request) (interface{}, error) {
	tr := &store.Translation{}

	err := json.NewDecoder(r.Body).Decode(tr)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

func PostTranslation(translations store.TranslationStore) StorePostFunc {
	return func(data interface{}) error {
		if tr, ok := data.(*store.Translation); !ok {
			return errors.New("not a translation")
		} else {
			return translations.Save(tr)
		}
	}
}