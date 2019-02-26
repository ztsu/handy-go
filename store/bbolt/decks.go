package bbolt

import (
	"encoding/json"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

const DecksBucketName = "Decks"
const UserDecksBucketName = "UserDecks"

type userDecks struct {
	UserID store.UUID   `json:"userId"`
	Decks  []store.UUID `json:"decks"`
}

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

		_, err = tx.CreateBucketIfNotExists([]byte(UserDecksBucketName))
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

func (ds *DecksBboltStore) Get(uuid store.UUID) (store.Deck, error) {
	deck := store.Deck{};

	return deck, ds.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DecksBucketName)).Get(uuid.MarshalBinary())

		return json.Unmarshal(b, &deck)
	})
}

func appendDeckToUserDecks(ud []store.UUID, deckID store.UUID) []store.UUID {
	for _, id :=range ud {
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

		err = tx.Bucket([]byte(DecksBucketName)).Put(deck.UUID.MarshalBinary(), b)
		if err != nil {
			return err
		}

		ud := userDecks{}
		b = tx.Bucket([]byte(UserDecksBucketName)).Get(deck.UserID.MarshalBinary())

		err = json.Unmarshal(b, &ud)
		if err != nil {
			return err
		}

		ud.Decks = appendDeckToUserDecks(ud.Decks, deck.UUID)

		udb, err := json.Marshal(ud)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(UserDecksBucketName)).Put(deck.UserID.MarshalBinary(), udb)
	})
}

func (ds *DecksBboltStore) Delete(deck store.Deck) error {
	return ds.db.Update(func (tx *bbolt.Tx) error {
		ud := userDecks{}
		b := tx.Bucket([]byte(UserDecksBucketName)).Get(deck.UserID.MarshalBinary())
		err := json.Unmarshal(b, &ud)
		if err != nil {
			return err
		}

		tmp := []store.UUID{}
		for _, id := range ud.Decks {
			if deck.UUID != id {
				tmp = append(tmp, id)
			}
		}

		ud.Decks = tmp

		udb, err := json.Marshal(ud)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte(UserDecksBucketName)).Put(deck.UserID.MarshalBinary(), udb)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(DecksBucketName)).Delete(deck.UUID.MarshalBinary())
	})
}

