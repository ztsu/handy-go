package bolt

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

const UsersBucketName = "Users"

type UserBoltStore struct {
	db *bbolt.DB
}

func NewUserBoltStore(db *bbolt.DB) (*UserBoltStore, error) {
	ts := &UserBoltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(UsersBucketName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	ts.db = db

	return ts, nil
}

func (us *UserBoltStore) Get(id uuid.UUID) (*store.User, error) {
	user := &store.User{}

	return user, us.db.View(func(tx *bbolt.Tx) error {
		key, err := id.MarshalBinary()
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(UsersBucketName)).Get(key)

		return json.Unmarshal(b, user)
	})
}

func (us *UserBoltStore) Add(u *store.User) error {
	return us.db.Update(func(tx *bbolt.Tx) error {
		key, err := u.ID.MarshalBinary()
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(UsersBucketName)).Get(key)
		if len(b) != 0 {
			return errors.New("user is exists")
		}

		value, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(UsersBucketName)).Put(key, value)
	})
}

func (us *UserBoltStore) Save(u *store.User) error {
	return us.db.Update(func(tx *bbolt.Tx) error {
		value, err := json.Marshal(u)
		if err != nil {
			return err
		}

		key, err := u.ID.MarshalBinary()
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(UsersBucketName)).Put(key, value)
	})
}
