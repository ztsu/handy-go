package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r.Handle("/metrics", promhttp.Handler())

	r.Group(func(r chi.Router) {

		r.Use(Prom)

		r.Route("/cards", func(r chi.Router) {
			r.Post("/", NewPostHandler(DecodeCard, PostCard(cards)))

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(idCtx)

				r.Get("/", NewGetHandler(getIDCtx, GetCard(cards)))

				r.Put("/", NewPutHandler(getIDCtx, DecodeCard, PutCard(cards)))

				r.Delete("/", NewDeleteHandler(getIDCtx, DeleteCard(cards)))
				r.Delete("/{DI}", NewDeleteHandler(getIDCtx, DeleteCard(cards)))
			})
		})

		r.Route("/deck-cards", func(r chi.Router) {
			r.Post("/", NewPostHandler(DecodeDeckCard, PostDeckCard(deckCards)))

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(idCtx)

				r.Get("/", NewGetHandler(getIDCtx, GetDeckCard(deckCards)))

				r.Put("/", NewPutHandler(getIDCtx, DecodeDeckCard, PutDeckCard(deckCards)))

				r.Delete("/", NewDeleteHandler(getIDCtx, DeleteDeckCard(deckCards)))
			})
		})

		r.Route("/decks", func(r chi.Router) {
			r.Post("/", NewPostHandler(DecodeDeck, PostDeck(decks)))

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(idCtx)

				r.Get("/", NewGetHandler(getIDCtx, GetDeck(decks)))

				r.Put("/", NewPutHandler(getIDCtx, DecodeDeck, PutDeck(decks)))

				r.Delete("/", NewDeleteHandler(getIDCtx, DeleteDeck(decks)))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", NewPostHandler(DecodeUser, PostUser(users)))

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(idCtx)

				r.Get("/", NewGetHandler(getIDCtx, GetUser(users)))

				r.Put("/", NewPutHandler(getIDCtx, DecodeUser, PutUser(users)))

				r.Delete("/", NewDeleteHandler(getIDCtx, DeleteUser(users)))
			})
		})
	})

	return r
}
