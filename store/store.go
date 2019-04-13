package store

import (
	"errors"
)

type Identity interface {
	Identity() string
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (user *User) Identity() string {
	return user.ID
}

type Card struct {
	ID     string `json:"id"`
	DeckID string `json:"deckId"`
	Type   string `json:"type"`
	Viewed uint64 `json:"viewed"`
	Opened uint64 `json:"opened"`
}

type Deck struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type Translation struct {
	ID          string `json:"id"`
	From        string `json:"from"`
	To          string `json:"to"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

type CardStore interface {
	Get(string) (Card, error)
	Save(Card) error
}

type DeckStore interface {
	Get(string) (Deck, error)
	Save(Deck) error
	Delete(Deck) error
}

type TranslationStore interface {
	Get(string) (Translation, error)
	Save(*Translation) error
	Delete(string) error
}

type UserStore interface {
	Add(*User) error
	Get(string) (*User, error)
	Save(*User) error
	Delete(string) error
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
