package bbolt

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
	"go.etcd.io/bbolt"
)

const DecksBucketName = "Decks"
const UserDecksBucketName = "UserDecks"

type userDecks struct {
	UserID uuid.UUID   `json:"userId"`
	Decks  []uuid.UUID `json:"decks"`
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

func (ds *DecksBboltStore) Get(id uuid.UUID) (store.Deck, error) {
	deck := store.Deck{};

	return deck, ds.db.View(func(tx *bbolt.Tx) error {
		key, err := id.MarshalBinary()
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(DecksBucketName)).Get(key)

		return json.Unmarshal(b, &deck)
	})
}

func appendDeckToUserDecks(ud []uuid.UUID, deckID uuid.UUID) []uuid.UUID {
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

		deckKey, err := deck.ID.MarshalBinary()
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte(DecksBucketName)).Put(deckKey, b)
		if err != nil {
			return err
		}

		userKey, err := deck.UserID.MarshalBinary()
		if err != nil {
			return err
		}

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
	return ds.db.Update(func (tx *bbolt.Tx) error {
		userKey, err := deck.UserID.MarshalBinary()
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(UserDecksBucketName)).Get(userKey)

		ud := userDecks{}
		err = json.Unmarshal(b, &ud)
		if err != nil {
			return err
		}

		tmp := []uuid.UUID{}
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

		deckKey, err := deck.ID.MarshalBinary()
		if err != nil {
			return err
		}

		return tx.Bucket([]byte(DecksBucketName)).Delete(deckKey)
	})
}

