package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
	"log"
	"net/http"
)

var (
	ErrCantParseID  = NewJsonError("can't parse id", http.StatusBadRequest)
	ErrUserNotFound = NewJsonError("user not found", http.StatusNotFound)
)

type DecodeFunc func(r *http.Request) (interface{}, error)

type StorePostFunc func(interface{}) error

type StorePutFunc func(interface{}) error

type StoreGetFunc func(uuid.UUID) (interface{}, error)

type StoreDeleteFunc func(uuid.UUID) error

type GetIDFromContextFunc func(ctx context.Context) string

func WriteJSON(w http.ResponseWriter, v interface{}, code int) {
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Panic(err)
	}
}

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

func NewGetHandler(getID GetIDFromContextFunc, get StoreGetFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(getID(r.Context()))
		if err != nil {
			ErrCantParseID.WriteTo(w)
			return
		}

		entity, err := get(uuid.UUID(id))
		if err != nil {
			ErrUserNotFound.WriteTo(w)
			return
		}

		WriteJSON(w, entity, 200)
	}
}

func NewDeleteHandler(getID GetIDFromContextFunc, delete StoreDeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(getID(r.Context()))
		if err != nil {
			ErrCantParseID.WriteTo(w)
			return
		}

		err = delete(id)
		if err != nil {
			ErrUserNotFound.WriteTo(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
