package main

import (
	"fmt"
	"github.com/go-chi/chi"
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

	translations, err := bbolt.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	decks, err := bbolt.NewDecksBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter();

	r.Route("/decks", func (r chi.Router) {
		r.Get(
			"/{ID}",
			store.NewGetFromStoreHandler(store.GetDeck(decks)),
		)
	})

	r.Route("/translations", func(r chi.Router) {

		r.Post("/", store.NewCreateTranslationHandler(translations))

		r.Get(
			"/{ID}",
			store.NewGetFromStoreHandler(store.GetTranslation(translations)),
		)

		r.Delete(
			"/{ID}",
			store.NewDeleteFromStoreHandler(store.DeleteTranslation(translations)),
		)
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}