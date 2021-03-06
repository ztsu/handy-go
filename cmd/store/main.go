package main

import (
	"database/sql"
	"fmt"
	store "github.com/ztsu/handy-go/store/http"
	"github.com/ztsu/handy-go/store/postgres"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
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

	r := store.NewRouter(cards, deckCards, decks, users)

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}
