package graphql

import (
	"context"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/ztsu/handy-go/store"
	"github.com/ztsu/handy-go/store/api"
	"net/http"
)

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

func (r *mutationResolver) Math(ctx context.Context) (*MathMutation, error) {
	return &MathMutation{}, nil
}