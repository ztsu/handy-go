package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	handy "github.com/ztsu/handy-go/store"
)

const (
	usersTableName = "users"

	usersPKeyConstraint = "users_pkey"
)

type UserStorePostgres struct {
	db *sql.DB
}

func NewUserStorePostgres(db *sql.DB) (*UserStorePostgres, error) {
	return &UserStorePostgres{db: db}, nil
}

func (store *UserStorePostgres) Add(user *handy.User) error {
	_, err := store.db.Exec("INSERT INTO users(id, email) VALUES($1, $2)", user.ID, user.Email)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Constraint == usersPKeyConstraint {
			return handy.ErrUserAlreadyExists
		}

		return err
	}

	return nil
}

func (store *UserStorePostgres) Get(id string) (*handy.User, error) {
	row := store.db.QueryRow("SELECT id, email FROM users WHERE id = $1", id)

	user, err := scan(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, handy.ErrUserNotFound
	}

	return user, err
}

func (store *UserStorePostgres) Save(user *handy.User) error {
	query := `UPDATE ` + usersTableName + ` SET email = $2 WHERE id = $1`
	res, err := store.db.Exec(query, user.ID, user.Email)
	if err != nil {
		return err
	}

	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if updated == 0 {
		return handy.ErrUserNotFound
	}

	return nil
}

func (store *UserStorePostgres) Delete(id string) error {
	query := `DELETE FROM ` + usersTableName + ` WHERE id = $1`
	res, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return handy.ErrUserNotFound
	}

	return nil
}

func scan(row *sql.Row) (*handy.User, error) {
	user := handy.User{}

	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
