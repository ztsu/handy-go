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
		return nil, ErrCantParseJson
	}

	return u, nil
}

func PostUser(users store.UserStore) StorePostFunc {
	return func(data interface{}) error {
		if user, ok := data.(*store.User); !ok {
			return errors.New("not a User")
		} else {
			return users.Add(user)
		}
	}
}

func GetUser(users store.UserStore) StoreGetFunc {
	return func(id string) (interface{}, error) { return users.Get(id) }
}

func PutUser(users store.UserStore) StorePutFunc {
	return func(data interface{}) error {
		if user, ok := data.(*store.User); !ok {
			return errors.New("not a User")
		} else {
			return users.Save(user)
		}
	}
}

func DeleteUser(users store.UserStore) StoreDeleteFunc {
	return func(id string) error { return users.Delete(id) }
}
