package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/internal/handy"
	"github.com/ztsu/handy-go/internal/store"
	"net/http"
)

func NewCreateDeckHandler(s *handy.UserService) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		u, err:= uuid.Parse(r.Header.Get("User-ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		userID := store.UUID(u)
		deck := store.Deck{UserID: userID}

		json.NewDecoder(r.Body).Decode(&deck)

		err = s.CreateDeck(userID, deck)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(deck)
	}
}

func NewDeleteDeckHandler(s *handy.UserService) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		u, err:= uuid.Parse(r.Header.Get("User-ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s.DeleteDeck(store.UUID(u), store.UUID(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}