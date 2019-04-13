package http

import (
	"bytes"
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"log"
	"net/http"
)

type JsonError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewJsonError(message string, code int) *JsonError {
	return &JsonError{Message: message, Code: code}
}

func (err *JsonError) Error() string {
	return err.Message
}

func (err *JsonError) WriteTo(w http.ResponseWriter) {
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.Code)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Print(err)
	}
}

var (
	ErrInternalServerError = NewJsonError("internal server error", http.StatusInternalServerError)
	ErrCantParseID         = NewJsonError("can't parse id", http.StatusBadRequest)
	ErrCantParseJson       = NewJsonError("can't parse json", http.StatusUnprocessableEntity)
	ErrIdentityMismatch    = NewJsonError("identity mismatch", http.StatusBadRequest)

	ErrUserNotFound      = NewJsonError("user not found", http.StatusNotFound)
	ErrUserAlreadyExists = NewJsonError("user already exists", http.StatusBadRequest)
)

var storeToJSONErrorMapping = map[error]*JsonError{
	store.ErrUserNotFound:      ErrUserNotFound,
	store.ErrUserAlreadyExists: ErrUserAlreadyExists,
}

func ConvertStoreErrorToJSONError(err error) *JsonError {
	if e, ok := err.(*JsonError); ok {
		return e
	}

	if e, ok := storeToJSONErrorMapping[err]; ok {
		return e
	}

	return ErrInternalServerError
}
