package handy

import "errors"

type Card struct {
	UUID     UUID   `json:"uuid"`
	DeckUUID UUID   `json:"deckUuid"`
	Type     string `json:"type"`
	Viewed   uint64 `json:"viewed"`
	Opened   uint64 `json:"opened"`
}

type CardStore interface {
	Get(UUID) (Card, error)
	Save(Card) error
}

type Deck struct {
	UUID        UUID   `json:"uuid"`
	UserID      UUID   `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type DeckStore interface {
	Get(UUID) (Deck, error)
	Save(Deck) error
	Delete(Deck) error
}

type Translation struct {
	UUID        UUID   `json:"uuid"`
	From        string `json:"from"`
	To          string `json:"to"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

type TranslationStore interface {
	Get(UUID) (Translation, error)
	Save(Translation) error
	Delete(UUID) error
}

type UserService struct {
	deckStore DeckStore
}

func NewUserService(ds DeckStore) *UserService {
	return &UserService{deckStore: ds}
}

func (s *UserService) CreateDeck(userID UUID, deck Deck) error {
	err := s.deckStore.Save(deck)

	return err
}

func (s *UserService) DeleteDeck(userID UUID, deckID UUID) error {
	deck, err := s.deckStore.Get(deckID)
	if err != nil {
		return err
	}

	if deck.UserID != userID {
		return errors.New("forbidden")
	}

	return s.deckStore.Delete(deck)
}