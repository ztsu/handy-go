package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http HttpClient
	basePath string
}

func NewClient(client HttpClient, basePath string) *Client {
	return &Client{http: client, basePath: basePath}
}

func (c *Client) AddUser(user *store.User) error {
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.basePath+"/users", b)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return store.ErrUserAlreadyExists
	default:
		return fmt.Errorf("Error: %d", resp.StatusCode)
	}
}