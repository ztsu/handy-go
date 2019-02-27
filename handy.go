package handy

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go/store"
)

type UserService struct {
	deckStore store.DeckStore
}

func NewUserService(ds store.DeckStore) *UserService {
	return &UserService{deckStore: ds}
}

func (s *UserService) CreateDeck(userID uuid.UUID, deck store.Deck) error {
	return s.deckStore.Save(deck)
}

func (s *UserService) DeleteDeck(userID uuid.UUID, deckID store.UUID) error {
	deck, err := s.deckStore.Get(deckID)
	if err != nil {
		return err
	}

	if deck.UserID != userID {
		return errors.New("forbidden")
	}

	return s.deckStore.Delete(deck)
}
