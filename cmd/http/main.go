package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ztsu/handy-go/pkg/handy"
	"net/http"
)

func DecksHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(handy.SampleDecks)
}

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(handy.SampleCards);
}

func main() {
	r := chi.NewRouter();

	r.Get("/decks", DecksHandler)
	r.Get("/cards", CardsHandler)

	addr := "0.0.0.0:8080"

	fmt.Printf("Starting server at %s\n", addr)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		fmt.Printf("Can't start server: %s\n", err)
	}
}