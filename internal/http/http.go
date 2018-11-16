package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/internal/store"
	"net/http"
)

func NewGetFormStoreHandler(s func (store.UUID) (interface{}, error)) func (http.ResponseWriter, *http.Request) {
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

func NewDeleteFormStoreHandler(s func (store.UUID) error) func (http.ResponseWriter, *http.Request) {
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
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func NewCreateTranslationHandler(repository store.TranslationStore) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
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