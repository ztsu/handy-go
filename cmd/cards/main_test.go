package main_test

import "testing"
import (
	handy "github.com/ztsu/handy-go"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	request := handy.Request{}

	response, err := handy.HandleRequest(nil, request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
}
