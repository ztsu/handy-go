// +build integration

package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCards(t *testing.T) {

	mux := newMux()

	u := bytes.NewBufferString(`{"id":"user-01","email":"user-01@example.org"}`)

	{
		req, _ := http.NewRequest("POST", "/users", u)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	}


	b := bytes.NewBufferString(`
{
	"id": "card-01",
	"userId":"user-01",
	"from":"Test deck",
	"to":"",
	"word":"",
	"translation":"",
	"ipa":""
}`)

	{
		req, _ := http.NewRequest("POST", "/cards", b)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	}

	{
		req, _ := http.NewRequest("DELETE", "/cards/card-01", nil)
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