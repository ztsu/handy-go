package http

import (
	"bytes"
	"encoding/json"
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
	WriteJsonError(w, err, err.Code)
}

func WriteJsonError(w http.ResponseWriter, err error, code int) {
	b := &bytes.Buffer{}

	if err := json.NewEncoder(b).Encode(err); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Panic(err)
	}
}
