package bolt

import (
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

const DecksBucketName = "Decks"
const UserDecksBucketName = "UserDecks"

type userDecks struct {
	UserID string   `json:"userId"`
	Decks  []string `json:"decks"`
}

type DecksBboltStore struct {
	db *bbolt.DB
}

func NewDecksBboltStore(db *bbolt.DB) (*DecksBboltStore, error) {
	ds := &DecksBboltStore{}

	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DecksBucketName))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(UserDecksBucketName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	ds.db = db

	return ds, nil
}

func (ds *DecksBboltStore) Get(id string) (store.Deck, error) {
	deck := store.Deck{}

	return deck, ds.db.View(func(tx *bbolt.Tx) error {
		key := []byte(id)

		b := tx.Bucket([]byte(DecksBucketName)).Get(key)

		return json.Unmarshal(b, &deck)
	})
}

func appendDeckToUserDecks(ud []string, deckID string) []string {
	for _, id := range ud {
		if id == deckID {
			return ud
		}
	}

	return append(ud, deckID)
}

func (ds *DecksBboltStore) Save(deck store.Deck) error {
	return ds.db.Update(func(tx *bbolt.Tx) error {
		b, err := json.Marshal(deck)
		if err != nil {
			return err
		}

		deckKey := []byte(deck.ID)

		err = tx.Bucket([]byte(DecksBucketName)).Put(deckKey, b)
		if err != nil {
			return err
		}

		userKey := []byte(deck.UserID)

		b = tx.Bucket([]byte(UserDecksBucketName)).Get(userKey)

		ud := userDecks{}

		err = json.Unmarshal(b, &ud)
		if err != nil {
			return err
		}

		ud.Decks = appendDeckToUserDecks(ud.Decks, deck.ID)

		udb, err := json.Marshal(ud)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(UserDecksBucketName)).Put(userKey, udb)
	})
}

func (ds *DecksBboltStore) Delete(deck store.Deck) error {
	return ds.db.Update(func(tx *bbolt.Tx) error {
		userKey := []byte(deck.UserID)

		b := tx.Bucket([]byte(UserDecksBucketName)).Get(userKey)

		ud := userDecks{}
		err := json.Unmarshal(b, &ud)
		if err != nil {
			return err
		}

		tmp := []string{}
		for _, id := range ud.Decks {
			if deck.ID != id {
				tmp = append(tmp, id)
			}
		}

		ud.Decks = tmp

		udb, err := json.Marshal(ud)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte(UserDecksBucketName)).Put(userKey, udb)
		if err != nil {
			return err
		}

		deckKey := []byte(deck.ID)

		return tx.Bucket([]byte(DecksBucketName)).Delete(deckKey)
	})
}
