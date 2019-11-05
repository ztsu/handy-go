// +build integration

package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecks(t *testing.T) {

	mux := newMux()

	u := bytes.NewBufferString(`{"id":"user-01","email":"user-01@example.org"}`)

	{
		req, _ := http.NewRequest("POST", "/users", u)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	}


	b := bytes.NewBufferString(`{"id":"deck-01","userId":"user-01","name":"Test deck"}`)

	{
		req, _ := http.NewRequest("POST", "/decks", b)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	}

	{
		req, _ := http.NewRequest("DELETE", "/decks/deck-01", nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	}

	{
		req, _ := http.NewRequest("DELETE", "/users/user-01", nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	}
}