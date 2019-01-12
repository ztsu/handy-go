package handy

import (
	"errors"
	"github.com/ztsu/handy-go/internal/store"
)

type UserService struct {
	deckStore store.DeckStore
}

func NewUserService(ds store.DeckStore) *UserService {
	return &UserService{deckStore: ds}
}

func (s *UserService) CreateDeck(userID store.UUID, deck store.Deck) error {
	err := s.deckStore.Save(deck)

	return err
}

func (s *UserService) DeleteDeck(userID store.UUID, deckID store.UUID) error {
	deck, err := s.deckStore.Get(deckID)
	if err != nil {
		return err
	}

	if deck.UserID != userID {
		return errors.New("forbidden")
	}

	return s.deckStore.Delete(deck)
}