package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ztsu/handy-go/internal/handy"
	httpHandy "github.com/ztsu/handy-go/internal/http"
	"github.com/ztsu/handy-go/internal/store"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
)

func main() {
	db, err := bbolt.Open("store.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ts, err := store.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	ds, err := store.NewDecksBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	dd := handy.NewUserService(ds)

	r := chi.NewRouter();

	r.Route("/api", func (router chi.Router) {
		router.Post("/decks", httpHandy.NewCreateDeckHandler(dd))

		router.Delete("/decks/{ID}", httpHandy.NewDeleteDeckHandler(dd))
	})

	r.Route("/int", func (r chi.Router) {
		r.Get(
			"/decks/{ID}",
			httpHandy.NewGetFormStoreHandler(func(id store.UUID) (interface{}, error) { return ds.Get(id) }),
		)

		r.Post("/translations", httpHandy.NewCreateTranslationHandler(ts))

		r.Get(
			"/translations/{ID}",
			httpHandy.NewGetFormStoreHandler(func(id store.UUID) (interface{}, error) { return ts.Get(id) }),
		)

		r.Delete(
			"/translations/{ID}",
			httpHandy.NewDeleteFormStoreHandler(func(id store.UUID) error { return ts.Delete(id) }),
		)
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}