package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	handy "github.com/ztsu/handy-go/store"
)

const (
	decksTableName = "decks"

	decksPkeyConstraint       = "decks_pkey"
	decksUserIdFkeyConstraint = "decks_userId_fkey"
)

type DeckStorePostgres struct {
	db *sql.DB
}

func NewDeckStorePostgres(db *sql.DB) (*DeckStorePostgres, error) {
	return &DeckStorePostgres{db: db}, nil
}

func (store *DeckStorePostgres) Add(deck *handy.Deck) error {
	query := `INSERT INTO ` + decksTableName + `(id, "userId", "name") VALUES($1, $2, $3)`
	_, err := store.db.Exec(query, deck.ID, deck.UserID, deck.Name)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Constraint == decksPkeyConstraint {
			return handy.ErrDeckAlreadyExists
		}
		if err, ok := err.(*pq.Error); ok && err.Constraint == decksUserIdFkeyConstraint {
			return handy.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (store *DeckStorePostgres) Get(id string) (*handy.Deck, error) {
	query := `SELECT id, "userId", "name" FROM ` + decksTableName + ` WHERE id = $1`
	row := store.db.QueryRow(query, id)
	println("qq")

	deck, err := scanDeck(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, handy.ErrDeckNotFound
	}
	println("qz")
	fmt.Printf("Err: %s\n", err)

	return deck, err
}

func (store *DeckStorePostgres) Save(deck *handy.Deck) error {
	query := `UPDATE ` + decksTableName + ` SET "userId" = $2, "name" = $3 WHERE id = $1`
	res, err := store.db.Exec(query, deck.ID, deck.UserID, deck.Name)
	if err != nil {
		return err
	}

	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if updated == 0 {
		return handy.ErrDeckNotFound
	}

	return nil
}

func (store *DeckStorePostgres) Delete(id string) error {
	query := `DELETE FROM ` + decksTableName + ` WHERE id = $1`
	res, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return handy.ErrDeckNotFound
	}

	return nil
}

func scanDeck(row *sql.Row) (*handy.Deck, error) {
	deck := handy.Deck{}

	err := row.Scan(&deck.ID, &deck.UserID, &deck.Name)
	if err != nil {
		return nil, err
	}

	return &deck, nil
}
