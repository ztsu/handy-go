//go:generate go run github.com/99designs/gqlgen

package graphql

import (
	"context"
	"github.com/matoous/go-nanoid"
	"github.com/ztsu/handy-go/store"
	"github.com/ztsu/handy-go/store/api"
	"net/http"
)

const UserID_TODO = "PLcB_sNkdStqiKMMfX4GAl"

type Resolver struct {
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Deck() DeckResolver {
	return &deckResolver{r}
}

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return "0", nil
}

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input RegisterUser) (*User, error) {
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
		return nil, err
	}

	return &User{ID: user.ID, Email: user.Email}, nil
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

type deckResolver struct {
	*Resolver
}

func (deckResolver) User(ctx context.Context, obj *Deck) (*User, error) {
	client := api.NewClient(http.DefaultClient, "http://localhost:8080")

	user, err := client.GetUser(obj.User.ID)
	if err != nil {
		return nil, err
	}

	return &User{ID: user.ID, Email: user.Email}, nil
}
