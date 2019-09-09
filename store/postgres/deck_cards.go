package postgres

import (
	"database/sql"
	"github.com/pkg/errors"
	handy "github.com/ztsu/handy-go/store"
)

const (
	deckCardsTableName = "deck_cards"
)

type DeckCardStorePostgres struct {
	db *sql.DB
}

func NewDeckCardStorePostgres(db *sql.DB) (*DeckCardStorePostgres, error) {
	return &DeckCardStorePostgres{db: db}, nil
}

func (store *DeckCardStorePostgres) Add(card *handy.DeckCard) error {
	return errors.New("not implemented yet")
}

func (store *DeckCardStorePostgres) Get(id string) (*handy.DeckCard, error) {
	return nil, errors.New("not implemented yet")
}

func (store *DeckCardStorePostgres) Save(card *handy.DeckCard) error {
	return errors.New("not implemented yet")
}

func (store *DeckCardStorePostgres) Delete(id string) error {
	return errors.New("not implemented yet")
}
