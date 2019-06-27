package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ztsu/handy-go/store/bolt"
	"github.com/ztsu/handy-go/store/dynamodb"
	store "github.com/ztsu/handy-go/store/http"
	"go.etcd.io/bbolt"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := bbolt.Open("store.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	translations, err := bolt.NewTranslationsBboltStore(db)
	if err != nil {
		log.Fatal(err)
	}

	//users, err := bolt.NewUserBoltStore(db)
	//if err != nil {
	//	log.Fatal(err)
	//}

	accessKeyID := os.Getenv("DYNAMODB_ACCESS_KEY_ID")
	secret := os.Getenv("DYNAMODB_ACCESS_KEY_SECRET")

	users, err := dynamodb.NewUserDynamoDBStore(accessKeyID, secret)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/translations", func(r chi.Router) {
		r.Post("/", store.NewPostHandler(store.DecodeTranslation, store.PostTranslation(translations)))

		r.Get(
			"/{ID}",
			store.NewGetHandler(store.GetID, store.GetTranslation(translations)),
		)

		r.Delete(
			"/{ID}",
			store.NewDeleteHandler(store.GetID, store.DeleteTranslation(translations)),
		)
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
