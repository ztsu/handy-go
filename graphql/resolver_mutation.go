package graphql

import (
	"context"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/ztsu/handy-go/store"
	"github.com/ztsu/handy-go/store/api"
	"net/http"
)

const UserID_TODO = "PLcB_sNkdStqiKMMfX4GAl"

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	output := &RegisterUserOutput{}

	client := api.NewClient(http.DefaultClient, "http://localhost:8080")

	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	user := &store.User{
		ID:    id,
		Email: input.Email,
	}

	err = client.AddUser(user)
	if err != nil {
		if e, ok := userErrorsToError(err).(*Error); ok {
			output.AddError(e)
			return output, nil
		}

		return output, err
	}

	output.User = &User{ID: user.ID, Email: user.Email}

	return output, nil
}


func (o *RegisterUserOutput) AddError(err *Error) {
	o.Ok = false
	o.Errors = append(o.Errors, err)
}

func userErrorsToError(err error) error {
	switch err {
	case store.ErrUserAlreadyExists:
		return ErrUserAlreadyExists;
	case store.ErrUserNotFound,
	store.ErrUserEmailNotProvided:
		return &Error{Message: err.Error()}
	}

	return err
}

func (r *mutationResolver) CreateDeck(ctx context.Context, input CreateDeck) (*Deck, error) {
	client := api.NewClient(http.DefaultClient, "http://localhost:8080")

	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	deck := &store.Deck{
		ID:     id,
		UserID: UserID_TODO,
		Name:   input.Name,
	}

	err = client.AddDeck(deck)
	if err != nil {
		return nil, err
	}

	return &Deck{
		ID:   deck.ID,
		Name: deck.Name,
		User: &User{ID: deck.UserID},
	}, nil
}