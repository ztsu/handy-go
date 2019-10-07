package http

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/ztsu/handy-go/store/postgres"
	"log"
	"os"
)

func newMux() *chi.Mux {
	pg, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	cards, err := postgres.NewCardStorePostgres(pg)
	if err != nil {
		log.Fatal(err)
	}

	deckCards, err := postgres.NewDeckCardStorePostgres(pg)
	if err != nil {
		log.Fatal(err)
	}

	decks, err := postgres.NewDeckStorePostgres(pg)
	if err != nil {
		log.Fatal(err)
	}

	users, err := postgres.NewUserStorePostgres(pg)
	if err != nil {
		log.Fatal(err)
	}

	r := NewRouter(cards, deckCards, decks, users)

	return r
}