package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go"
	"github.com/ztsu/handy-go/http/middleware"
	"net/http"
)

func NewCreateDeckHandler(s *handy.UserService) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())
		deck := handy.Deck{UserID: userID}

		err := json.NewDecoder(r.Body).Decode(&deck)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Error: %s", err)
		}

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

func NewDeleteDeckHandler(s *handy.UserService) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		userID := handy.UUID(middleware.GetUserID(r.Context()))

		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s.DeleteDeck(userID, handy.UUID(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
