package bolt

import (
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

var (
	UsersBucketName      = []byte("Users")
	UsersEmailBucketName = []byte("UserEmails")
)

type UserBoltStore struct {
	db *bbolt.DB
}

func NewUserBoltStore(db *bbolt.DB) (*UserBoltStore, error) {
	ts := &UserBoltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(UsersBucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(UsersEmailBucketName)
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

func (us *UserBoltStore) Get(id string) (*store.User, error) {
	user := &store.User{}

	return user, us.db.View(func(tx *bbolt.Tx) error {
		key := []byte(id)

		b := tx.Bucket(UsersBucketName).Get(key)

		if len(b) == 0 {
			return store.ErrUserNotFound
		}

		return json.Unmarshal(b, user)
	})
}

func (us *UserBoltStore) Add(user *store.User) error {
	return us.db.Update(func(tx *bbolt.Tx) error {
		err := user.Validate()
		if err != nil {
			return err
		}

		users := tx.Bucket(UsersBucketName)
		emails := tx.Bucket(UsersEmailBucketName)

		key := []byte(user.ID)

		if b := users.Get(key); len(b) != 0 {
			return store.ErrUserAlreadyExists
		}

		if len(user.Email) == 0 {
			return store.ErrUserEmailNotProvided
		}

		if b := emails.Get([]byte(user.Email)); len(b) != 0 {
			return store.ErrUserAlreadyExists
		}

		if err := emails.Put([]byte(user.Email), key); err != nil {
			return err
		}

		value, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return users.Put(key, value)
	})
}

func (us *UserBoltStore) Save(user *store.User) error {
	return us.db.Update(func(tx *bbolt.Tx) error {
		err := user.Validate()
		if err != nil {
			return err
		}

		key := []byte(user.ID)

		existed, err := us.Get(user.ID)
		if err != nil {
			return err
		}

		if user.Email != existed.Email {
			err := tx.Bucket(UsersEmailBucketName).Delete([]byte(existed.Email))
			if err != nil {
				return err
			}

			if err := tx.Bucket(UsersEmailBucketName).Put([]byte(user.Email), key); err != nil {
				return err
			}
		}

		value, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return tx.Bucket(UsersBucketName).Put(key, value)
	})
}

func (us *UserBoltStore) Delete(id string) error {
	return us.db.Update(func(tx *bbolt.Tx) error {
		key := []byte(id)

		if len(tx.Bucket(UsersBucketName).Get(key)) == 0 {
			return store.ErrUserNotFound
		}

		return tx.Bucket(UsersBucketName).Delete(key)
	})
}
