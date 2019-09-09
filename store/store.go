package store

import (
	"errors"
	"github.com/go-playground/validator"
)

type Identity interface {
	Identity() string
}

type User struct {
	ID    string `json:"id"    validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (user *User) Identity() string {
	return user.ID
}

func (user *User) Validate() error {
	err := validator.New().Struct(user)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		return ErrUserUnprocessable
	}

	return nil
}

type UserStore interface {
	Add(*User) error
	Get(string) (*User, error)
	Save(*User) error
	Delete(string) error
}

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserUnprocessable    = errors.New("user is unprocessable")
	ErrUserEmailNotProvided = errors.New("email not provided")
)

type Deck struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

func (deck *Deck) Identity() string {
	return deck.ID
}

type DeckStore interface {
	Add(*Deck) error
	Get(string) (*Deck, error)
	Save(*Deck) error
	Delete(string) error
}

var (
	ErrDeckAlreadyExists = errors.New("deck already exists")
	ErrDeckNotFound      = errors.New("deck not found")
)

type Card struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	From        string `json:"from"`
	To          string `json:"to"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

type CardStore interface {
	Add(card *Card) error
	Get(string) (*Card, error)
	Save(*Card) error
	Delete(string) error
}

var (
	ErrCardAlreadyExists = errors.New("card already exists")
	ErrCardNotFound      = errors.New("card not found")
)

type DeckCard struct {
	ID     string `json:"id"`
	DeckID string `json:"deckId"`
	CardID string `json:"cardId"`
	Views  uint64 `json:"views"`
	Turns  uint64 `json:"Turns"`
}

type DeckCardStore interface {
	Add(*DeckCard) error
	Get(string) (*DeckCard, error)
	Save(*DeckCard) error
	Delete(string) error
}

var (
	ErrDeckCardAlreadyExists = errors.New("card already exists in deck")
	ErrDeckCardNotFound      = errors.New("card not found in deck")
)
