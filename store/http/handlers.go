package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

type StoreGetFunc func(uuid2 store.UUID) (interface{}, error)

type StoreDeleteFunc func(store.UUID) error

//func NewCreateDeckHandler(s *handy.UserService) http.HandlerFunc {
//	return func (w http.ResponseWriter, r *http.Request) {
//		userID := middleware.GetUserID(r.Context())
//		deck := store.Deck{UserID: userID}
//
//		err := json.NewDecoder(r.Body).Decode(&deck)
//		if err != nil {
//			w.WriteHeader(http.StatusUnprocessableEntity)
//			fmt.Fprintf(w, "Error: %s", err)
//		}
//
//		err = s.CreateDeck(userID, deck)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
//			return
//		}
//
//		w.WriteHeader(http.StatusCreated)
//		json.NewEncoder(w).Encode(deck)
//	}
//}

func NewGetFromStoreHandler(s StoreGetFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		entity, err := s(store.UUID(id))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		json.NewEncoder(w).Encode(entity)
	}
}

func NewDeleteFromStoreHandler(s StoreDeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s(store.UUID(id))
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetDeck(decks store.DeckStore) StoreGetFunc {
	return func(id store.UUID) (interface{}, error) { return decks.Get(id) }
}

func GetTranslation(translations store.TranslationStore) StoreGetFunc {
	return func(id store.UUID) (interface{}, error) { return translations.Get(id) }
}

func DeleteTranslation(translations store.TranslationStore) StoreDeleteFunc {
	return func(id store.UUID) error { return translations.Delete(id) }
}

func NewCreateTranslationHandler(repository store.TranslationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tr store.Translation

		err := json.NewDecoder(r.Body).Decode(&tr)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = repository.Save(tr)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		json.NewEncoder(w).Encode(tr)
	}
}

