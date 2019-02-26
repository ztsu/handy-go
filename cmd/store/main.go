package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ztsu/handy-go/store/bbolt"
	store "github.com/ztsu/handy-go/store/http"
	"log"
	"net/http"
)

func main() {
	db, err := bbolt.Open("store.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()


	decks, err := bbolt.NewDecksBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	translations, err := bbolt.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	users, err := bbolt.NewUserBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter();

	r.Use(middleware.Logger)

	r.Route("/decks", func (r chi.Router) {
		r.Get(
			"/{ID}",
			store.NewGetHandler(store.GetDeck(decks)),
		)
	})

	r.Route("/translations", func(r chi.Router) {

		r.Post("/", store.NewPostHandler(store.DecodeTranslation, store.PostTranslation(translations)))

		r.Get(
			"/{ID}",
			store.NewGetHandler(store.GetTranslation(translations)),
		)

		r.Delete(
			"/{ID}",
			store.NewDeleteHandler(store.DeleteTranslation(translations)),
		)
	})


	r.Route("/users", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeUser, store.PostUser(users)))

		r.Put("/{ID}", store.NewPutHandler(store.DecodeUser, store.PutUser(users)))
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}