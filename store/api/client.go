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
	http     HttpClient
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

func (c *Client) GetUser(id string) (*store.User, error) {
	req, err := http.NewRequest("GET", c.basePath+"/users/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		user := new(store.User)
		err := json.NewDecoder(resp.Body).Decode(user)
		return user, err
	case http.StatusNotFound:
		return nil, store.ErrUserNotFound
	default:
		return nil, fmt.Errorf("Error: %d", resp.StatusCode)
	}
}

func (c *Client) AddDeck(deck *store.Deck) error {
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(deck)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.basePath+"/decks", b)
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
	default:
		return fmt.Errorf("Error: %d", resp.StatusCode)
	}
}
