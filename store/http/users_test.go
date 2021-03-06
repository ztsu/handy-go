// +build integration

package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsers(t *testing.T) {
	mux := newMux()

	b := bytes.NewBufferString(`{"id":"test-01","email":"test-01@example.org"}`)

	{
		req, _ := http.NewRequest("POST", "/users", b)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	}

	{
		req, _ := http.NewRequest("GET", "/users/test-01", nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	}

	{
		req, _ := http.NewRequest("DELETE", "/users/test-01", nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	}

	{
		req, _ := http.NewRequest("GET", "/users/test-01", nil)
		resp := httptest.NewRecorder()

		mux.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	}
}
