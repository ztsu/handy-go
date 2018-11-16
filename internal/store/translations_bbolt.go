package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
)

const TranslationsBucketName = "Translations"

type TranslationsBboltStore struct {
	db *bbolt.DB
}

func NewTranslationsBboltStore(db *bbolt.DB) (*TranslationsBboltStore, error) {
	store := &TranslationsBboltStore{}

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

	store.db = db

	return store, nil
}

func (repository *TranslationsBboltStore) Get(uuid UUID) (Translation, error) {
	tr := Translation{};

	return tr, repository.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(TranslationsBucketName)).Get(uuid.MarshalBinary())

		return json.Unmarshal(b, &tr)
	})
}

func (repository *TranslationsBboltStore) Delete(uuid UUID) error {
	return repository.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(TranslationsBucketName)).Delete(uuid.MarshalBinary())
	})
}

func (repository *TranslationsBboltStore) Save(tr Translation) error {
	return repository.db.Update(func(tx *bbolt.Tx) error {
		b, err := json.Marshal(tr)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(TranslationsBucketName)).Put(tr.UUID.MarshalBinary(), b)
	})
}