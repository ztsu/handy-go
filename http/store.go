package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go"
	"net/http"
)

func NewGetFormStoreHandler(s func (handy.UUID) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		entity, err := s(handy.UUID(id))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		json.NewEncoder(w).Encode(entity)
	}
}

func NewDeleteFormStoreHandler(s func (handy.UUID) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s(handy.UUID(id))
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
