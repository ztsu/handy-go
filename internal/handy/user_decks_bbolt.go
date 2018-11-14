package handy

import (
	"encoding/json"
	"go.etcd.io/bbolt"
)

const UserDecksBucketName = "UserDecks"

type UserDecksBboltStore struct {
	db *bbolt.DB
}

func NewUserDecksBboltStore(db *bbolt.DB) (*UserDecksBboltStore, error) {
	store := &UserDecksBboltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(UserDecksBucketName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	store.db = db

	return store, nil
}

func (store *UserDecksBboltStore) Get(userID UUID) (UserDecks, error) {
	ud := UserDecks{};

	return ud, store.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(UserDecksBucketName)).Get(userID.MarshalBinary())

		return json.Unmarshal(b, &ud)
	})
}

func (store *UserDecksBboltStore) Delete(userID UUID) error {
	return store.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(UserDecksBucketName)).Delete(userID.MarshalBinary())
	})
}

func (store *UserDecksBboltStore) Save(ud UserDecks) error {
	return store.db.Update(func(tx *bbolt.Tx) error {
		b, err := json.Marshal(ud)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(UserDecksBucketName)).Put(ud.UserID.MarshalBinary(), b)
	})
}