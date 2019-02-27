package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

type DecodeFunc func(r *http.Request) (interface{}, error)

type StorePostFunc func(interface{}) error

type StorePutFunc func(interface{}) error

type StoreGetFunc func(uuid.UUID) (interface{}, error)

type StoreDeleteFunc func(uuid.UUID) error

func NewPostHandler(m DecodeFunc, s StorePostFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := m(r)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s(data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func NewPutHandler(m DecodeFunc, s StorePutFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))

		data, err := m(r)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		type WithID struct {
			ID uuid.UUID
		}

		if idty, ok := data.(store.Identity); !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if idty.Identity() != id {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "identity mismatch")
			return
		}

		err = s(data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func NewGetHandler(s StoreGetFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		entity, err := s(uuid.UUID(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		json.NewEncoder(w).Encode(entity)
	}
}

func NewDeleteHandler(s StoreDeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "ID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = s(uuid.UUID(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
