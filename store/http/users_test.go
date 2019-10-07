// +build integration

package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	b := bytes.NewBufferString(`{"id":"test-01","email":"test-01@example.org"}`)

	req, _ := http.NewRequest("POST", "/users", b)
	resp := httptest.NewRecorder()

	newMux().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "OK response is expected")
}