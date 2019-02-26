package store

import (
	original "github.com/google/uuid"
)

type Identity interface {
	Identity() original.UUID
}

type User struct {
	ID    original.UUID `json:"id"`
	Email string        `json:"email"`
}

func (user *User) Identity() original.UUID {
	return user.ID
}

type Card struct {
	ID     UUID   `json:"id"`
	DeckID UUID   `json:"deckId"`
	Type   string `json:"type"`
	Viewed uint64 `json:"viewed"`
	Opened uint64 `json:"opened"`
}

type Deck struct {
	ID          UUID   `json:"id"`
	UserID      UUID   `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type Translation struct {
	ID          UUID   `json:"uuid"`
	From        string `json:"from"`
	To          string `json:"to"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

type CardStore interface {
	Get(UUID) (Card, error)
	Save(Card) error
}

type DeckStore interface {
	Get(UUID) (Deck, error)
	Save(Deck) error
	Delete(Deck) error
}

type TranslationStore interface {
	Get(UUID) (Translation, error)
	Save(*Translation) error
	Delete(UUID) error
}

type UserStore interface {
	Add(*User) error
	Get(original.UUID) (*User, error)
	Save(*User) error
}
