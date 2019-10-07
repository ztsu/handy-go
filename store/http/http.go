package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ztsu/handy-go/store"
)

func NewRouter(
	cards store.CardStore,
	deckCards store.DeckCardStore,
	decks store.DeckStore,
	users store.UserStore,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/cards", func(r chi.Router) {
		r.Post("/", NewPostHandler(DecodeCard, PostCard(cards)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(QueryStringID("ID"))

			r.Get("/", NewGetHandler(GetID, GetCard(cards)))

			r.Put("/", NewPutHandler(GetID, DecodeCard, PutCard(cards)))

			r.Delete("/", NewDeleteHandler(GetID, DeleteCard(cards)))
		})
	})

	r.Route("/deck-cards", func(r chi.Router) {
		r.Post("/", NewPostHandler(DecodeDeckCard, PostDeckCard(deckCards)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(QueryStringID("ID"))

			r.Get("/", NewGetHandler(GetID, GetDeckCard(deckCards)))

			r.Put("/", NewPutHandler(GetID, DecodeDeckCard, PutDeckCard(deckCards)))

			r.Delete("/", NewDeleteHandler(GetID, DeleteDeckCard(deckCards)))
		})
	})

	r.Route("/decks", func(r chi.Router) {
		r.Post("/", NewPostHandler(DecodeDeck, PostDeck(decks)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(QueryStringID("ID"))

			r.Get("/", NewGetHandler(GetID, GetDeck(decks)))

			r.Put("/", NewPutHandler(GetID, DecodeDeck, PutDeck(decks)))

			r.Delete("/", NewDeleteHandler(GetID, DeleteDeck(decks)))
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", NewPostHandler(DecodeUser, PostUser(users)))

		r.Route("/{ID}", func(r chi.Router) {
			r.Use(QueryStringID("ID"))

			r.Get("/", NewGetHandler(GetID, GetUser(users)))

			r.Put("/", NewPutHandler(GetID, DecodeUser, PutUser(users)))

			r.Delete("/", NewDeleteHandler(GetID, DeleteUser(users)))
		})
	})

	return r
}