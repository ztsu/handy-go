package main_test

import "testing"
import (
	"github.com/stretchr/testify/assert"
	handy "github.com/ztsu/handy-go"
)

func TestHandleRequest(t *testing.T) {
	request := handy.Request{}

	response, err := handy.HandleRequest(nil, request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
}
