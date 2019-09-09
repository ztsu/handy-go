package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	handy "github.com/ztsu/handy-go/store"
)

const (
	cardsTableName = "cards"

	cardsPkeyConstraint       = "cards_pkey"
	cardsUserIdFkeyConstraint = "cards_userId_fkey"
)

type CardStorePostgres struct {
	db *sql.DB
}

func NewCardStorePostgres(db *sql.DB) (*CardStorePostgres, error) {
	return &CardStorePostgres{db: db}, nil
}

func (store *CardStorePostgres) Add(card *handy.Card) error {
	query := `INSERT INTO ` + cardsTableName +
		`(id, "userId", "from", "to", "word", "translation", "ipa") VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := store.db.Exec(query, card.ID, card.UserID, card.From, card.To, card.Word, card.Translation, card.IPA)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Constraint == cardsPkeyConstraint {
			return handy.ErrCardAlreadyExists
		}
		if err, ok := err.(*pq.Error); ok && err.Constraint == cardsUserIdFkeyConstraint {
			return handy.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (store *CardStorePostgres) Get(id string) (*handy.Card, error) {
	return nil, errors.New("not implemented yet")
}

func (store *CardStorePostgres) Save(card *handy.Card) error {
	return errors.New("not implemented yet")
}

func (store *CardStorePostgres) Delete(id string) error {
	query := `DELETE FROM ` + cardsTableName + ` WHERE id = $1`
	res, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return handy.ErrCardNotFound
	}

	return nil
}
