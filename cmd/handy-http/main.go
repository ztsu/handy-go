package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ztsu/handy-go"
	handyBbolt "github.com/ztsu/handy-go/bbolt"
	handyHttp "github.com/ztsu/handy-go/http"
	"github.com/ztsu/handy-go/http/middleware"
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

	ts, err := handyBbolt.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	ds, err := handyBbolt.NewDecksBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	dd := handy.NewUserService(ds)

	r := chi.NewRouter();

	r.Route("/api", func (router chi.Router) {
		router.With(middleware.Authenticator).Post("/decks", handyHttp.NewCreateDeckHandler(dd))

		router.Delete("/decks/{ID}", handyHttp.NewDeleteDeckHandler(dd))
	})

	r.Route("/int", func (r chi.Router) {
		r.Get(
			"/decks/{ID}",
			handyHttp.NewGetFormStoreHandler(func(id handy.UUID) (interface{}, error) { return ds.Get(id) }),
		)

		r.Post("/translations", handyHttp.NewCreateTranslationHandler(ts))

		r.Get(
			"/translations/{ID}",
			handyHttp.NewGetFormStoreHandler(func(id handy.UUID) (interface{}, error) { return ts.Get(id) }),
		)

		r.Delete(
			"/translations/{ID}",
			handyHttp.NewDeleteFormStoreHandler(func(id handy.UUID) error { return ts.Delete(id) }),
		)
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}