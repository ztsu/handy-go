package handy

import (
	"encoding/json"
	"go.etcd.io/bbolt"
)

const DecksBucketName = "Decks"

type DecksBboltStore struct {
	db *bbolt.DB
}

func NewDecksBboltStore(db *bbolt.DB) (*DecksBboltStore, error) {
	store := &DecksBboltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DecksBucketName))
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

func (repository *DecksBboltStore) Get(uuid UUID) (Deck, error) {
	deck := Deck{};

	return deck, repository.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DecksBucketName)).Get(uuid.MarshalBinary())

		return json.Unmarshal(b, &deck)
	})
}

func (repository *DecksBboltStore) Delete(uuid UUID) error {
	return repository.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(DecksBucketName)).Delete(uuid.MarshalBinary())
	})
}

func (repository *DecksBboltStore) Save(tr Deck) error {
	return repository.db.Update(func(tx *bbolt.Tx) error {
		b, err := json.Marshal(tr)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(DecksBucketName)).Put(tr.UUID.MarshalBinary(), b)
	})
}