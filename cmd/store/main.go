package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	store "github.com/ztsu/handy-go/store/http"
	"github.com/ztsu/handy-go/store/postgres"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := bbolt.Open("store.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	pg, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	//

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

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/cards", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeCard, store.PostCard(cards)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(store.QueryStringID("ID"))

			r.Get("/", store.NewGetHandler(store.GetID, store.GetCard(cards)))

			r.Put("/", store.NewPutHandler(store.GetID, store.DecodeCard, store.PutCard(cards)))

			r.Delete("/", store.NewDeleteHandler(store.GetID, store.DeleteCard(cards)))
		})
	})

	r.Route("/deck-cards", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeDeckCard, store.PostDeckCard(deckCards)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(store.QueryStringID("ID"))

			r.Get("/", store.NewGetHandler(store.GetID, store.GetDeckCard(deckCards)))

			r.Put("/", store.NewPutHandler(store.GetID, store.DecodeDeckCard, store.PutDeckCard(deckCards)))

			r.Delete("/", store.NewDeleteHandler(store.GetID, store.DeleteDeckCard(deckCards)))
		})
	})

	r.Route("/decks", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeDeck, store.PostDeck(decks)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(store.QueryStringID("ID"))

			r.Get("/", store.NewGetHandler(store.GetID, store.GetDeck(decks)))

			r.Put("/", store.NewPutHandler(store.GetID, store.DecodeDeck, store.PutDeck(decks)))

			r.Delete("/", store.NewDeleteHandler(store.GetID, store.DeleteDeck(decks)))
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeUser, store.PostUser(users)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(store.QueryStringID("ID"))

			r.Get("/", store.NewGetHandler(store.GetID, store.GetUser(users)))

			r.Put("/", store.NewPutHandler(store.GetID, store.DecodeUser, store.PutUser(users)))

			r.Delete("/", store.NewDeleteHandler(store.GetID, store.DeleteUser(users)))
		})
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}
