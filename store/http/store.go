package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Panic(err)
	}
}

type DecodeFunc func(r *http.Request) (interface{}, error)

type StorePostFunc func(interface{}) error

func NewPostHandler(decode DecodeFunc, fn StorePostFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := decode(r)
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		err = fn(data)
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		writeJSON(w, data, http.StatusCreated)
	}
}

type GetIDFromContextFunc func(ctx context.Context) string

type StoreGetFunc func(string) (interface{}, error)

func NewGetHandler(getID GetIDFromContextFunc, fn StoreGetFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := fn(getID(r.Context()))
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		writeJSON(w, res, http.StatusOK)
	}
}

type StorePutFunc func(interface{}) error

func NewPutHandler(getID GetIDFromContextFunc, decode DecodeFunc, fn StorePutFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := decode(r)
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		if idty, ok := data.(store.Identity); !ok {
			ErrInternalServerError.WriteTo(w)
			return
		} else if idty.Identity() != getID(r.Context()) {
			ErrIdentityMismatch.WriteTo(w)
			return
		}

		err = fn(data)
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		writeJSON(w, data, http.StatusOK)
	}
}

type StoreDeleteFunc func(string) error

func NewDeleteHandler(getID GetIDFromContextFunc, fn StoreDeleteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(getID(r.Context()))
		if err != nil {
			ConvertStoreErrorToJSONError(err).WriteTo(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
