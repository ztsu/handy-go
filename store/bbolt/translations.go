package bbolt

import (
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

const TranslationsBucketName = "Translations"

type TranslationsBboltStore struct {
	db *bbolt.DB
}

func NewTranslationsBboltStore(db *bbolt.DB) (*TranslationsBboltStore, error) {
	ts := &TranslationsBboltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TranslationsBucketName))
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

func (ts *TranslationsBboltStore) Get(uuid store.UUID) (store.Translation, error) {
	tr := store.Translation{};

	return tr, ts.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TranslationsBucketName)).Get(uuid.MarshalBinary())

		return json.Unmarshal(b, &tr)
	})
}

func (ts *TranslationsBboltStore) Delete(uuid store.UUID) error {
	return ts.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(TranslationsBucketName)).Delete(uuid.MarshalBinary())
	})
}

func (ts *TranslationsBboltStore) Save(tr store.Translation) error {
	return ts.db.Update(func(tx *bbolt.Tx) error {
		b, err := json.Marshal(tr)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(TranslationsBucketName)).Put(tr.UUID.MarshalBinary(), b)
	})
}