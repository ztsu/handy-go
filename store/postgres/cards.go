package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	handy "github.com/ztsu/handy-go/store"
)

const (
	cardsTableName = "cards"

	cardsPKeyConstraint       = "cards_pkey"
	cardsUserIdFKeyConstraint = "cards_userId_fkey"
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
		if err, ok := err.(*pq.Error); ok && err.Constraint == cardsPKeyConstraint {
			return handy.ErrCardAlreadyExists
		}
		if err, ok := err.(*pq.Error); ok && err.Constraint == cardsUserIdFKeyConstraint {
			return handy.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (store *CardStorePostgres) Get(id string) (*handy.Card, error) {
	query := `SELECT id, "userId", "from", "to", "word", "translation", "ipa" FROM ` + cardsTableName + ` WHERE id = $1`
	row := store.db.QueryRow(query, id)

	card, err := scanCard(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, handy.ErrCardNotFound
	}

	return card, err
}

func (store *CardStorePostgres) Save(card *handy.Card) error {
	query := `UPDATE ` + cardsTableName + `
SET "userId" = $2, "from" = $3, "to" = $4, "word" = $5, "translation" = $6, "ipa" = $7
WHERE id = $1`
	res, err := store.db.Exec(query, card.ID, card.UserID, card.From, card.To, card.Word, card.Translation, card.IPA)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Constraint == cardsUserIdFKeyConstraint {
			return handy.ErrUserNotFound
		}
		return err
	}

	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if updated == 0 {
		return handy.ErrCardNotFound
	}

	return nil
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

func scanCard(row *sql.Row) (*handy.Card, error) {
	card := &handy.Card{}

	err := row.Scan(&card.ID, &card.UserID, &card.From, &card.To, &card.Word, &card.Translation, &card.IPA)
	if err != nil {
		return nil, err
	}

	return card, nil
}
