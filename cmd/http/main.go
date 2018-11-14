package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ztsu/handy-go/internal/handy"
	httpHandy "github.com/ztsu/handy-go/internal/http"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
)

func main() {
	db, err := bbolt.Open("handy.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	tr, err := handy.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	dr, err := handy.NewDecksBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter();

	r.Route("/api", func (r chi.Router) {
		r.Post("/decks", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.Header.Get("User-ID")))
		})
	})

	r.Route("/int", func (r chi.Router) {
		r.Get(
			"/decks/{ID}",
			httpHandy.NewGetFormStoreHandler(func(id handy.UUID) (interface{}, error) { return dr.Get(id) }),
		)

		r.Post("/translations", httpHandy.NewCreateTranslationHandler(tr))

		r.Get(
			"/translations/{ID}",
			httpHandy.NewGetFormStoreHandler(func(id handy.UUID) (interface{}, error) { return tr.Get(id) }),
		)

		r.Delete(
			"/translations/{ID}",
			httpHandy.NewDeleteFormStoreHandler(func(id handy.UUID) error { return tr.Delete(id) }),
		)
	})

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}