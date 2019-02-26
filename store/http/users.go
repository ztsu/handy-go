package http

import (
	"encoding/json"
	"errors"
	"github.com/ztsu/handy-go/store"
	"net/http"
)

func DecodeUser(r *http.Request) (interface{}, error) {
	u := &store.User{}

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func PostUser(users store.UserStore) StorePostFunc {
	return func(data interface{}) error {
		if user, ok := data.(*store.User); !ok {
			return errors.New("not a user")
		} else {
			return users.Add(user)
		}
	}
}

func PutUser(users store.UserStore) StorePutFunc {
	return func(data interface{}) error {
		if user, ok := data.(*store.User); !ok {
			return errors.New("not a user")
		} else {
			return users.Save(user)
		}
	}
}